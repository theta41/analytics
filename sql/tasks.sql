CREATE TABLE public.tasks (
	id SERIAL NOT NULL CONSTRAINT tasks_pk PRIMARY KEY,
	object_id int NOT NULL,
	created_at int NOT NULL,
	finished_at int NOT NULL,
	is_done BOOLEAN NOT NULL
);