-- name: RetrieveChirp :one
SELECT * FROM chirps
WHERE id = $1;