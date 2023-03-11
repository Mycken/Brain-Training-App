insert into bta.public.voc_part_sp (descrip) values ('noun'),('verb'),('adjective'),('adverbe'),('predicate');

INSERT INTO bta.public.voc_test (descrip) VALUES ('Shulte'),('Arithmetic'),('Memorize Words');

INSERT INTO bta.public.users (username, email, password) VALUES ('Majkl','majkl@gmail.com','123');
INSERT INTO bta.public.users (username, email, password) VALUES ('Bob','bob@gmail.com','456');

INSERT INTO bta.public.res_test (user_id, test_id, date_test , result_inter) VALUES (1, 1, '2023-02-15','00:01:15')
