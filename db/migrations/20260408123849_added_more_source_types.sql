-- +goose Up
ALTER TYPE SourceType ADD VALUE 'Youtube';
ALTER TYPE SourceType ADD VALUE 'Other';

-- +goose Down
ALTER TYPE SourceType REMOVE VALUE 'Youtube';
ALTER TYPE SourceType REMOVE VALUE 'Other';
