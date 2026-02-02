package util

import "github.com/google/uuid"

func StringToUUID(s string) (uuid.UUID, error) {
    id, err := uuid.Parse(s)
    if err != nil {
        return uuid.Nil, err
    }
    return id, nil
}

func UUIDToString(id uuid.UUID) string {
    return id.String()
}