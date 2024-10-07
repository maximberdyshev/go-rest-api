INSERT INTO music_groups ("name") VALUES ('music_group_1');
INSERT INTO music_groups ("name") VALUES ('music_group_2');

INSERT INTO songs ("name", group_id, release_date, "text", "link")
VALUES ('song_1', 1, '01.01.1990', '{"kuplet_1", "kuplet_2", "kuplet_3", "kuplet_4"}', 'http://song_1_example.com');
INSERT INTO songs ("name", group_id, release_date, "text", "link")
VALUES ('song_2', 1, '10.01.1988', '{"kuplet_1", "kuplet_2", "kuplet_3", "kuplet_4"}', 'http://song_2_example.com');
INSERT INTO songs ("name", group_id, release_date, "text", "link")
VALUES ('song_3', 2, '12.12.2000', '{"kuplet_1", "kuplet_2", "kuplet_3"}', 'http://song_3_example.com');
INSERT INTO songs ("name", group_id, release_date, "text", "link")
VALUES ('song_4', 1, '07.09.2009', '{"kuplet_1", "kuplet_2", "kuplet_3", "kuplet_4"}', 'http://song_4_example.com');
INSERT INTO songs ("name", group_id, release_date, "text", "link")
VALUES ('song_5', 2, '11.04.1998', '{"kuplet_1", "kuplet_2"}', 'http://song_5_example.com');
