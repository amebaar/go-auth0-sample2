INSERT INTO user(cid,name) VALUES ((SELECT id FROM company WHERE name='A company'), 'A user 1');
INSERT INTO user(cid,name) VALUES ((SELECT id FROM company WHERE name='A company'), 'A user 2');
INSERT INTO user(cid,name) VALUES ((SELECT id FROM company WHERE name='A company'), 'A user 3');
INSERT INTO user(cid,name) VALUES ((SELECT id FROM company WHERE name='B company'), 'B user 1');
INSERT INTO user(cid,name) VALUES ((SELECT id FROM company WHERE name='B company'), 'B user 2');
INSERT INTO user(cid,name) VALUES ((SELECT id FROM company WHERE name='B company'), 'B user 3');