package storage

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/iamwavecut/ct-mend/internal/config"
	"github.com/iamwavecut/ct-mend/tools"
)

type MongoDB struct {
	ctx      context.Context
	conn     *mongo.Client
	clients  *mongo.Collection
	projects *mongo.Collection
	counters *mongo.Collection
	cancel   context.CancelFunc
	timeout  time.Duration
}

func (m *MongoDB) Init(ctx context.Context, connAddr string) error {
	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	m.ctx = ctx

	timeout := config.DefaultTimeout // Default
	cfg := ctx.Value(config.Key{})
	if cfg, ok := cfg.(*config.Config); ok && cfg != nil {
		timeout = cfg.GracefulTimeout
	}
	m.timeout = timeout

	opts := options.Client().ApplyURI(connAddr).SetConnectTimeout(timeout)
	client, err := mongo.Connect(ctx, opts)
	if !tools.Try(err) {
		return errors.Wrap(err, "mongo connect failed")
	}

	m.conn = client
	err = client.Ping(ctx, nil)
	if !tools.Try(err) {
		return errors.Wrap(err, "mongo ping failed")
	}
	m.clients = client.Database("mend").Collection("clients")
	m.projects = client.Database("mend").Collection("projects")
	m.counters = client.Database("mend").Collection("counters")
	return nil
}

func (m *MongoDB) Shutdown() {
	tools.Try(m.conn.Disconnect(m.ctx), true)
}

func (m *MongoDB) SelectClients() ([]*Client, error) {
	//goland:noinspection ALL
	res := []*Client{}
	ctx := m.getCtx()
	cur, err := m.clients.Find(ctx, bson.D{})
	if !tools.Try(err) {
		return nil, passNotFound(err)
	}

	for cur.Next(ctx) {
		client := &Client{}
		err = cur.Decode(client)
		if !tools.Try(err) {
			return nil, err
		}
		res = append(res, client)
	}
	err = cur.Err()
	if !tools.Try(err) {
		return res, passNotFound(err)
	}
	err = cur.Close(ctx)
	if !tools.Try(err) {
		return res, passNotFound(err)
	}

	if len(res) == 0 {
		return res, ErrNotFound{}
	}
	return res, nil
}

func (m *MongoDB) GetClient(ID int) (*Client, error) {
	var res *Client
	mres := m.clients.FindOne(m.getCtx(), bson.M{"id": ID})
	if mres == nil {
		return nil, ErrNotFound{}
	}
	err := mres.Err()
	if !tools.Try(err) {
		return nil, passNotFound(err)
	}
	err = mres.Decode(&res)
	if !tools.Try(err) {
		return nil, err
	}
	return res, nil
}

func (m *MongoDB) UpsertClient(client *Client) (*Client, error) {
	if client == nil {
		return nil, ErrNilEntity{}
	}

	if client.ID == nil {
		newID := m.newClientID()
		client.ID = &newID
	}
	ctx := m.getCtx()
	res, err := m.clients.UpdateOne(
		ctx,
		bson.M{"id": client.ID},
		bson.D{{Key: "$set", Value: client}},
		options.Update().SetUpsert(true),
	)
	if !tools.Try(err) {
		return nil, err
	}
	if res.UpsertedCount+res.ModifiedCount == 0 {
		return nil, ErrNotFound{}
	}
	newClient, err := m.GetClient(*client.ID)
	if !tools.Try(err) {
		return nil, err
	}
	return newClient, nil
}

func (m *MongoDB) DeleteClient(ID int) error {
	dres, err := m.clients.DeleteOne(m.getCtx(), bson.M{"id": ID})
	if !tools.Try(err) {
		return passNotFound(err)
	}
	if dres.DeletedCount == 0 {
		return ErrNotFound{}
	}
	return nil
}

func (m *MongoDB) SelectProjects() ([]*Project, error) {
	//goland:noinspection ALL
	res := []*Project{}
	ctx := m.getCtx()
	cur, err := m.projects.Find(ctx, bson.D{})
	if !tools.Try(err) {
		return nil, passNotFound(err)
	}

	for cur.Next(ctx) {
		project := &Project{}
		err = cur.Decode(project)
		if !tools.Try(err) {
			return nil, err
		}
		res = append(res, project)
	}
	err = cur.Err()
	if !tools.Try(err) {
		return res, passNotFound(err)
	}
	err = cur.Close(ctx)
	if !tools.Try(err) {
		return res, passNotFound(err)
	}

	if len(res) == 0 {
		return res, ErrNotFound{}
	}
	return res, nil
}

func (m *MongoDB) GetProject(ID int) (*Project, error) {
	var res *Project
	mres := m.projects.FindOne(m.getCtx(), bson.M{"id": ID})
	if mres == nil {
		return nil, ErrNotFound{}
	}
	err := mres.Err()
	if !tools.Try(err) {
		return nil, passNotFound(err)
	}
	err = mres.Decode(&res)
	if !tools.Try(err) {
		return nil, err
	}
	return res, nil
}

func (m *MongoDB) SelectProjectsOfClient(ID int) ([]*Project, error) {
	//goland:noinspection ALL
	res := []*Project{}
	ctx := m.getCtx()
	cur, err := m.projects.Find(ctx, bson.M{"client_id": ID})
	if !tools.Try(err) {
		return nil, passNotFound(err)
	}

	for cur.Next(ctx) {
		project := &Project{}
		err = cur.Decode(project)
		if !tools.Try(err) {
			return nil, err
		}
		res = append(res, project)
	}
	err = cur.Err()
	if !tools.Try(err) {
		return res, passNotFound(err)
	}
	err = cur.Close(ctx)
	if !tools.Try(err) {
		return res, passNotFound(err)
	}

	if len(res) == 0 {
		return res, ErrNotFound{}
	}
	return res, nil
}

func (m *MongoDB) UpsertProject(project *Project) (*Project, error) {
	if project == nil {
		return nil, ErrNilEntity{}
	}

	if project.ID == nil {
		newID := m.newProjectID()
		project.ID = &newID
	}
	ctx := m.getCtx()
	res, err := m.projects.UpdateOne(
		ctx,
		bson.M{"id": project.ID},
		bson.D{{Key: "$set", Value: project}},
		options.Update().SetUpsert(true),
	)
	if !tools.Try(err) {
		return nil, err
	}
	if res.UpsertedCount+res.ModifiedCount == 0 {
		return nil, ErrNotFound{}
	}
	newProject, err := m.GetProject(*project.ID)
	if !tools.Try(err) {
		return nil, err
	}
	return newProject, nil
}

func (m *MongoDB) DeleteProject(ID int) error {
	dres, err := m.projects.DeleteOne(m.getCtx(), bson.M{"id": ID})
	if !tools.Try(err) {
		return passNotFound(err)
	}
	if dres.DeletedCount == 0 {
		return ErrNotFound{}
	}
	return nil
}

func (m *MongoDB) getCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	time.AfterFunc(m.timeout+(1*time.Second), func() {
		cancel()
	})
	return ctx
}

func (m *MongoDB) newClientID() int {
	return m.newID("clients")
}

func (m *MongoDB) newProjectID() int {
	return m.newID("projects")
}

func (m *MongoDB) newID(key string) int {
	res := m.counters.FindOneAndUpdate(
		m.getCtx(),
		bson.D{{Key: "_id", Value: key}},
		bson.D{
			{Key: "$inc", Value: bson.D{{Key: "seq", Value: 1}}},
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)
	if !tools.Try(res.Err()) {
		if res.Err() == mongo.ErrNoDocuments {
			return 1
		} else {
			tools.Must(res.Err())
		}
	}

	ID := struct {
		_   string `bson:"_id"`
		Seq int    `bson:"seq"`
	}{}
	err := res.Decode(&ID)
	tools.Must(err)
	return ID.Seq
}
