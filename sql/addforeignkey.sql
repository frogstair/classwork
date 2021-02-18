ALTER TABLE requests
DROP CONSTRAINT IF EXISTS fk_assignments_requests;
ALTER TABLE requests
ADD CONSTRAINT fk_assignments_requests 
FOREIGN KEY (assignment_id) 
REFERENCES assignments (id);

ALTER TABLE assignment_files 
DROP CONSTRAINT IF EXISTS fk_assignments_files;
ALTER TABLE assignment_files 
ADD CONSTRAINT fk_assignments_files 
FOREIGN KEY (assignment_id) 
REFERENCES assignments (id);

ALTER TABLE request_uploads
DROP CONSTRAINT IF EXISTS fk_request_uploads;
ALTER TABLE request_uploads
ADD CONSTRAINT fk_request_uploads
FOREIGN KEY (request_id) 
REFERENCES requests (id);

ALTER TABLE assignments 
DROP CONSTRAINT IF EXISTS fk_subjects_assignments;
ALTER TABLE assignments 
ADD CONSTRAINT fk_subjects_assignments
FOREIGN KEY (subject_id) 
REFERENCES subjects (id);

ALTER TABLE subjects 
DROP CONSTRAINT IF EXISTS fk_teachers_subjects;
ALTER TABLE subjects 
ADD CONSTRAINT fk_teachers_subjects
FOREIGN KEY (teacher_id)
REFERENCES users (id);

ALTER TABLE assignments_completed
DROP CONSTRAINT IF EXISTS fk_users_assignments_completed;
ALTER TABLE assignments_completed
ADD CONSTRAINT fk_users_assignments_completed
FOREIGN KEY (user_id)
REFERENCES users (id);

ALTER TABLE assignments_completed
DROP CONSTRAINT IF EXISTS fk_assignments_assignments_completed;
ALTER TABLE assignments_completed
ADD CONSTRAINT fk_assignments_assignments_completed
FOREIGN KEY (assignment_id)
REFERENCES assignments (id);

ALTER TABLE school_students
DROP CONSTRAINT IF EXISTS fk_users_school;
ALTER TABLE school_students
ADD CONSTRAINT fk_users_school
FOREIGN KEY (user_id)
REFERENCES users (id);

ALTER TABLE school_students
DROP CONSTRAINT IF EXISTS fk_users_school;
ALTER TABLE school_students
ADD CONSTRAINT fk_users_school
FOREIGN KEY (school_id)
REFERENCES schools (id);

ALTER TABLE school_teachers
DROP CONSTRAINT IF EXISTS fk_users_school;
ALTER TABLE school_students
ADD CONSTRAINT fk_users_school
FOREIGN KEY (user_id)
REFERENCES users (id);

ALTER TABLE school_teachers
DROP CONSTRAINT IF EXISTS fk_users_school;
ALTER TABLE school_teachers
ADD CONSTRAINT fk_users_school
FOREIGN KEY (school_id)
REFERENCES schools (id);

ALTER TABLE school_subjects
DROP CONSTRAINT IF EXISTS fk_subject_school;
ALTER TABLE school_subjects
ADD CONSTRAINT fk_subject_school
FOREIGN KEY (subject_id)
REFERENCES subjects (id);

ALTER TABLE school_subjects
DROP CONSTRAINT IF EXISTS fk_subject_school;
ALTER TABLE school_subjects
ADD CONSTRAINT fk_subject_school
FOREIGN KEY (school_id)
REFERENCES schools (id);

ALTER TABLE subject_students
DROP CONSTRAINT IF EXISTS fk_user_subject;
ALTER TABLE subject_students
ADD CONSTRAINT fk_user_subject
FOREIGN KEY (subject_id)
REFERENCES subjects (id);

ALTER TABLE subject_students
DROP CONSTRAINT IF EXISTS fk_user_subject;
ALTER TABLE subject_students
ADD CONSTRAINT fk_user_subject
FOREIGN KEY (user_id)
REFERENCES users (id);