
create table if not exists snippets (
    id integer not null primary key auto_increment,
    title varchar(100) not null,
    content text not null,
    created datetime not null,
    expires datetime not null
);

create index idx_snippets_created on snippets(created);

INSERT INTO snippets (title, content, created, expires) VALUES (
                                                                   'An old silent pond',
                                                                   'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
                                                                   UTC_TIMESTAMP(),
                                                                   DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
                                                               );

INSERT INTO snippets (title, content, created, expires) VALUES (
                                                                   'Over the wintry forest',
                                                                   'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
                                                                   UTC_TIMESTAMP(),
                                                                   DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
                                                               );

INSERT INTO snippets (title, content, created, expires) VALUES (
                                                                   'First autumn morning',
                                                                   'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
                                                                   UTC_TIMESTAMP(),
                                                                   DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
                                                               );