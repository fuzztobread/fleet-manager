package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"fleet-manager/internal/vehicle"
)

// VehicleStore implements vehicle.Store using the shared DB.
type VehicleStore struct {
	db *DB
}

func NewVehicleStore(db *DB) *VehicleStore {
	return &VehicleStore{db: db}
}

func (s *VehicleStore) Create(ctx context.Context, v *vehicle.Vehicle) error {
	query := `
	INSERT INTO vehicles (id, name, status, capacity)
	VALUES (?, ?, ?, ?)`

	_, err := s.db.conn.ExecContext(ctx, query, v.ID, v.Name, v.Status, v.Capacity)
	if err != nil {
		return fmt.Errorf("create vehicle: %w", err)
	}

	return nil
}

// GetByID fetches a single vehicle by its primary key.
func (s *VehicleStore) GetByID(ctx context.Context, id string) (*vehicle.Vehicle, error) {
	query := `
	SELECT id, name, status, capacity
	FROM vehicles
	WHERE id = ?`

	var v vehicle.Vehicle
	err := s.db.conn.QueryRowContext(ctx, query, id).Scan(
		&v.ID,
		&v.Name,
		&v.Status,
		&v.Capacity,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, vehicle.ErrNotFound
		}
		return nil, fmt.Errorf("get vehicle: %w", err)
	}

	return &v, nil
}

func (s *VehicleStore) UpdateStatus(ctx context.Context, id string, status vehicle.Status) error {
	query := `
	UPDATE vehicles
	SET    status = ?
	WHERE  id = ?`

	result, err := s.db.conn.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update vehicle status: %w", err)
	}

	// If 0 rows were updated the ID doesn't exist — return the domain error.
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return vehicle.ErrNotFound
	}

	return nil
}

// List returns all vehicles — no filtering for now.
func (s *VehicleStore) List(ctx context.Context) ([]*vehicle.Vehicle, error) {
	query := `
	SELECT id, name, status, capacity
	FROM vehicles
	ORDER BY name ASC`

	rows, err := s.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list vehicles: %w", err)
	}
	// rows.Close() releases the DB connection back to the pool.
	defer rows.Close()
	vehicles := make([]*vehicle.Vehicle, 0)
	for rows.Next() {
		var v vehicle.Vehicle
		if err := rows.Scan(&v.ID, &v.Name, &v.Status, &v.Capacity); err != nil {
			return nil, fmt.Errorf("scan vehicle: %w", err)
		}
		vehicles = append(vehicles, &v)
	}

	// rows.Err() catches any error that happened during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate vehicles: %w", err)
	}

	return vehicles, nil
}
