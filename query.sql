-- name: SelectProfiles :many
SELECT * FROM "profile";

-- name: PhoneExisted :one
SELECT EXISTS(SELECT 1 FROM "profile" WHERE "last_name" = $1 LIMIT 1);

-- name: CreateProfile :one
INSERT INTO "profile" (
  "first_name", "last_name"
) VALUES (
  $1, $2
) RETURNING "id";

-- name: CreateUser :exec
INSERT INTO "users" (
"id", "phone_number"
) VALUES (
  $1 , $2 
);


-- name: SetOTP :exec
UPDATE "users"
  set "otp" = $2,
  "otp_expiration_time" = current_timestamp + (interval '1 minute')
WHERE "phone_number" = $1;


-- name: IsValidUserOTP :one
SELECT EXISTS(SELECT "otp_expiration_time" FROM "users" WHERE "otp" = $1 AND  phone_number = $2  LIMIT 1) AS valid;

-- name: IsOTPExpired :one
SELECT EXISTS(SELECT "otp_expiration_time" FROM "users" WHERE "otp" = $1 AND  phone_number = $2  AND "otp_expiration_time" > current_timestamp LIMIT 1) AS valid;
