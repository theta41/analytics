CREATE TABLE public.letters (
    id SERIAL NOT NULL CONSTRAINT letters_pk PRIMARY KEY,
    task_id INT NOT null CONSTRAINT letters_tasks_id_fk REFERENCES tasks ON UPDATE CASCADE ON DELETE CASCADE,
	email VARCHAR(320) NOT NULL,
    sent_at int NOT NULL,
	finished_at int NOT NULL,
	is_accepted BOOLEAN NOT NULL
);