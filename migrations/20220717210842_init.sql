-- +goose Up
-- +goose StatementBegin
	CREATE SCHEMA IF NOT EXISTS task;

	CREATE TABLE IF NOT EXISTS task.task_event_types(
		id serial PRIMARY KEY,
		task_type VARCHAR(20)
	);
	
	INSERT INTO task.task_event_types (task_type) VALUES ('created');
    INSERT INTO task.task_event_types (task_type) VALUES ('approved');
    INSERT INTO task.task_event_types (task_type) VALUES ('rejected');
    INSERT INTO task.task_event_types (task_type) VALUES ('send_mail');

	CREATE TABLE IF NOT EXISTS task.users(
		id serial PRIMARY KEY,
		email VARCHAR(50)
	);
	
	CREATE TABLE IF NOT EXISTS task.tasks(
		id serial PRIMARY KEY,
		author_id serial REFERENCES task.users (id),
		descr VARCHAR(80),
		body VARCHAR(800),
		finished BOOLEAN
		--status serial REFERENCES task.task_event_types (id)
	);
	
	CREATE TABLE IF NOT EXISTS task.task_approvers(
		id serial PRIMARY KEY,
		task_id serial REFERENCES task.tasks (id),
		approver_id serial REFERENCES task.users (id)
	);
	
	CREATE TABLE IF NOT EXISTS task.task_events(
		id serial PRIMARY KEY,
		task_id serial REFERENCES task.tasks (id),
		user_id serial REFERENCES task.users (id),
		event_type_id INTEGER,
		event_time timestamp
	);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
	DROP TABLE IF EXISTS task.task_events; 
	
	DROP TABLE IF EXISTS task.task_approvers;
	
	DROP TABLE IF EXISTS task.tasks;
	
	DROP TABLE IF EXISTS task.users;
	
	DROP TABLE IF EXISTS task.task_event_types;

-- +goose StatementEnd
