import React, { useState } from 'react';
import {
  Box, Container, Header, FormField,
  Input, Button, SpaceBetween, Flashbar
} from '@cloudscape-design/components';
import { useNavigate } from 'react-router-dom';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState([]);
  const navigate = useNavigate();

  const handleLogin = () => {
    if (!username || !password) {
      setErrors([{ type: 'error', content: 'Please enter both username and password.', id: 'error-1' }]);
    } else {
      setErrors([]);
      navigate('/dashboard');
    }
  };

  return (
    <Box padding="xxl" display="flex" justifyContent="center">
      <Container header={<Header variant="h1">Sign in</Header>} variant="stacked" fitHeight>
        <SpaceBetween size="l">
          {errors.length > 0 && <Flashbar items={errors} />}
          <FormField label="Username">
            <Input
              value={username}
              onChange={({ detail }) => setUsername(detail.value)}
              placeholder="Enter your username"
            />
          </FormField>
          <FormField label="Password">
            <Input
              type="password"
              value={password}
              onChange={({ detail }) => setPassword(detail.value)}
              placeholder="Enter your password"
            />
          </FormField>
          <Button variant="primary" onClick={handleLogin}>Sign in</Button>
        </SpaceBetween>
      </Container>
    </Box>
  );
}
