const request = require('supertest');
const app = require('../src/index');

describe('User Service', () => {
  describe('GET /health', () => {
    it('should return healthy status', async () => {
      const res = await request(app).get('/health');
      expect(res.statusCode).toBe(200);
      expect(res.body.status).toBe('healthy');
      expect(res.body.service).toBe('user-service');
    });
  });

  describe('GET /users', () => {
    it('should return list of users', async () => {
      const res = await request(app).get('/users');
      expect(res.statusCode).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
      expect(res.body.length).toBeGreaterThan(0);
    });
  });

  describe('GET /users/:id', () => {
    it('should return user by id', async () => {
      const res = await request(app).get('/users/1');
      expect(res.statusCode).toBe(200);
      expect(res.body.id).toBe(1);
      expect(res.body.name).toBe('John Doe');
    });

    it('should return 404 for non-existent user', async () => {
      const res = await request(app).get('/users/999');
      expect(res.statusCode).toBe(404);
    });
  });

  describe('POST /users', () => {
    it('should create a new user', async () => {
      const res = await request(app)
        .post('/users')
        .send({ name: 'Test User', email: 'test@example.com' });
      expect(res.statusCode).toBe(201);
      expect(res.body.name).toBe('Test User');
    });

    it('should return 400 if name is missing', async () => {
      const res = await request(app)
        .post('/users')
        .send({ email: 'test@example.com' });
      expect(res.statusCode).toBe(400);
    });
  });
});
