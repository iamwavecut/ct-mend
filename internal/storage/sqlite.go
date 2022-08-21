package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Driver for SQLite3 init

	"github.com/iamwavecut/ct-mend/tools"
)

type (
	SQLite struct {
		conn *sqlx.DB
		ctx  context.Context
	}

	flatClient struct {
		Client
		ClientSettings
	}
)

func (c *flatClient) Inflate() *Client {
	return &Client{
		ID:       c.ID,
		Name:     c.Name,
		Settings: c.ClientSettings,
	}
}

func (s *SQLite) Init(ctx context.Context, path string) error {
	db, err := sqlx.Connect("sqlite3", path)
	if !tools.Try(err) {
		return err
	}
	s.conn = db
	s.ctx = ctx
	return nil
}

func (s *SQLite) SelectClients() ([]*Client, error) {
	//goland:noinspection ALL
	proxy := []*flatClient{}
	err := s.conn.SelectContext(s.ctx, &proxy, `
			select id, name, code_scan_interval from clients;
		`)
	if !tools.Try(err) {
		return nil, passNotFound(err)
	}
	var res []*Client
	for _, flat := range proxy {
		res = append(res, flat.Inflate())
	}
	return res, nil
}

func (s *SQLite) GetClient(ID int) (*Client, error) {
	proxy := flatClient{}
	err := s.conn.GetContext(s.ctx, &proxy, "select id, name, code_scan_interval from clients where id=?;", ID)
	if !tools.Try(err) {
		return nil, passNotFound(err)
	}
	return proxy.Inflate(), nil
}

func (s *SQLite) UpsertClient(client *Client) (*Client, error) {
	if client == nil {
		return nil, ErrNilEntity{}
	}
	res := &flatClient{}
	err := s.conn.GetContext(s.ctx, res, `
		insert into clients (id, name, code_scan_interval) 
		values (?,?,?) 
		on conflict(id) do update set
		    name=excluded.name,
			code_scan_interval=excluded.code_scan_interval
		where excluded.id=id
		returning id, name, code_scan_interval;
	`, client.ID, client.Name, client.Settings.CodeScanInterval)
	if !tools.Try(err) {
		return nil, err
	}
	return res.Inflate(), nil
}

func (s *SQLite) DeleteClient(ID int) error {
	res, err := s.conn.ExecContext(s.ctx, "delete from clients where id = ?;", ID)
	if !tools.Try(err) {
		return err
	}
	n, err := res.RowsAffected()
	if !tools.Try(err) {
		return err
	}
	if n == 0 {
		return ErrNotFound{}
	}
	return nil
}

func (s *SQLite) SelectProjects() ([]*Project, error) {
	//goland:noinspection ALL
	res := []*Project{}
	err := s.conn.SelectContext(s.ctx, &res, `
				select id, client_id, name from projects;
			`)
	if !tools.Try(err) {
		return nil, ErrNotFound{}
	}
	return res, nil
}

func (s *SQLite) GetProject(ID int) (*Project, error) {
	res := &Project{}
	err := s.conn.GetContext(s.ctx, res, "select id, client_id, name from projects where id=?;", ID)
	if !tools.Try(err) {
		return nil, passNotFound(err)
	}
	return res, nil
}

func (s *SQLite) SelectProjectsOfClient(ID int) ([]*Project, error) {
	var res []*Project
	err := s.conn.SelectContext(s.ctx, &res, `select id, client_id, name from projects where client_id=?;`, ID)
	if !tools.Try(err) {
		return nil, err
	}
	return res, nil
}

func (s *SQLite) UpsertProject(project *Project) (*Project, error) {
	if project == nil {
		return nil, ErrNilEntity{}
	}
	res := &Project{}
	err := s.conn.GetContext(s.ctx, res, `
		insert into projects (id, client_id, name) 
		values (?,?,?) 
		on conflict(id) do update set
		    name=excluded.name,
			client_id=excluded.client_id
		where excluded.id=id
		returning id, client_id, name;
	`, project.ID, project.ClientID, project.Name)
	if !tools.Try(err) {
		return nil, err
	}
	return res, nil
}

func (s *SQLite) DeleteProject(ID int) error {
	res, err := s.conn.ExecContext(s.ctx, "delete from projects where id = ?;", ID)
	if !tools.Try(err) {
		return err
	}
	n, err := res.RowsAffected()
	if !tools.Try(err) {
		return err
	}
	if n == 0 {
		return ErrNotFound{}
	}
	return nil
}
