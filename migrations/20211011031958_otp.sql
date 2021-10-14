-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.otps (
    id uuid NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id uuid,
    code varchar(4),
    used boolean
);
ALTER TABLE ONLY public.otps ADD CONSTRAINT otps_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.otps ADD CONSTRAINT fk_otps_user FOREIGN KEY (user_id) REFERENCES public.users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE otps;
-- +goose StatementEnd
