import React, { useState } from 'react';
import { Container, Row, Col, Form, Button, Card } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

export default function SignIn() {

  const NAMESPACE = 'http://localhost:8080/api';

  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

  const [errorMessage, setErrorMessage] = useState('');
  const navigate = useNavigate();

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log('Form Data:', formData);
  
    try {
      const response = await axios.post(`${NAMESPACE}/signin`, formData, {
        headers: {
          'Content-Type': 'application/json',
        },
      });
  
      if (response.status === 200) {
        setErrorMessage('');
        console.log('Login successful!');
      } else {
        console.error('Login failed. Invalid credentials.');
        setErrorMessage('Login failed. Invalid credentials.');
      }
    } catch (error) {
      console.error('Login failed. Invalid credentials.');
      setErrorMessage('Login failed. Invalid credentials.');
    }
  
    // Reset the form after submission
    setFormData({
      email: '',
      password: '',
    });
  };

  return (
    <div>
    <Container className="mt-5">
      <Row className="justify-content-md-center">
        <Col xs={12} md={6}>
          <Card>
            <Card.Header as="h5">Login</Card.Header>
            <Card.Body>
              {errorMessage && <div className="alert alert-danger">{errorMessage}</div>}
              <Form onSubmit={handleSubmit}>
                <Form.Group controlId="formEmail" className="mb-4">
                  <Form.Label>Email Address</Form.Label>
                  <Form.Control
                    type="email"
                    placeholder="Enter your email"
                    name="email"
                    value={formData.email}
                    onChange={handleInputChange}
                    required
                  />
                </Form.Group>

                <Form.Group controlId="formPassword" className="mb-4">
                  <Form.Label>Password</Form.Label>
                  <Form.Control
                    type="password"
                    placeholder="Enter your password"
                    name="password"
                    value={formData.password}
                    onChange={handleInputChange}
                    required
                  />
                </Form.Group>

                <Button variant="primary" type="submit">
                  Login
                </Button>
              </Form>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
    </div>
  );
}
