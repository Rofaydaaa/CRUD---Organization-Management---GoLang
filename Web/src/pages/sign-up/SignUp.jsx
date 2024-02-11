import React, { useState } from 'react';
import { Container, Row, Col, Form, Button, Card, Alert } from 'react-bootstrap';
import axios from 'axios';

export default function SignUp() {

  const NAMESPACE = 'http://localhost:8080/api';

  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: ''
  });

  const [errorMessage, setErrorMessage] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prevFormData => ({
      ...prevFormData,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const response = await axios.post(`${NAMESPACE}/signup`, formData, {
        headers: {
          'Content-Type': 'application/json',
        },
      });
    
      if (response.status === 200) {
        setErrorMessage('');
        setSuccessMessage('Registration successful!');
        console.log('Registration successful!');
      } else {
        setErrorMessage('Registration failed. Please try again.');
        console.error('Registration failed. Please try again.');
      }
    } catch (error) {
      setErrorMessage(`Error during registration: ${error.message}`);
      console.error('Error during registration:', error.message);
      console.error('Registration response:', error.response.data);
    }

    // Reset the form after submission
    setFormData({
      name: '',
      email: '',
      password: ''
    });
  };

  return (
    <div className='mb-5'>
      <Container className="mt-5">
        <Row className="justify-content-md-center">
          <Col xs={12} md={6}>
            <Card>
              <Card.Header as="h5">Sign Up</Card.Header>
              <Card.Body>
                {errorMessage && <Alert variant="danger">{errorMessage}</Alert>}
                {successMessage && <Alert variant="success">{successMessage}</Alert>}
                <Form onSubmit={handleSubmit}>
                  <Form.Group controlId="formName" className="mb-4">
                    <Form.Label>Name</Form.Label>
                    <Form.Control
                      type="text"
                      placeholder="Enter your name"
                      name="name"
                      value={formData.name}
                      onChange={handleInputChange}
                      required
                    />
                  </Form.Group>

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
                    Sign Up
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
