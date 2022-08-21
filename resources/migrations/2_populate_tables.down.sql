delete from clients;

delete from projects;

update `sqlite_sequence`
set `seq` = 0
where `name` in ('clients', 'projects');