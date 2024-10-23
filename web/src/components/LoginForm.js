import React, { useEffect, useState } from 'react';
import {
  Button,
  Divider,
  Form,
  Grid,
  Header,
  Image,
  Message,
  Modal,
  Segment,
} from 'semantic-ui-react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useUser } from '../context/User';
import { API, showError, showSuccess } from '../helpers';


// Mock API for login
const mockLogin = async (username, password) => {
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 1000));

  // Check credentials (you can add more users here)
  const validUsers = [
    { username: 'admin', password: '123', role: '1000' },
    { username: 'user', password: '123', role: '1' },
  ];

  const user = validUsers.find(u => u.username === username && u.password === password);

  if (user) {
    return {
      success: true,
      message: 'Login successful',
      data: {
        token: 'mock-jwt-token',
        user: { id: 1, username: user.username, role: user.role }
      }
    };
  } else {
    return {
      success: false,
      message: 'Invalid username or password',
    };
  }
};

const LoginForm = () => {
  const [inputs, setInputs] = useState({
    username: '',
    password: '',
    wechat_verification_code: '',
  });
  const [searchParams, setSearchParams] = useSearchParams();
  const [submitted, setSubmitted] = useState(false);
  const { username, password } = inputs;
  const { login } = useUser();
  let navigate = useNavigate();

  const [status, setStatus] = useState({});

  useEffect(() => {
    if (searchParams.get("expired")) {
      showError('未登录或登录已过期，请重新登录！');
    }
    let status = localStorage.getItem('status');
    if (status) {
      status = JSON.parse(status);
      setStatus(status);
    }
  }, []);

  const [showWeChatLoginModal, setShowWeChatLoginModal] = useState(false);

  const onGitHubOAuthClicked = () => {
    window.open(
      `https://github.com/login/oauth/authorize?client_id=${status.github_client_id}&scope=user:email`
    );
  };

  const onWeChatLoginClicked = () => {
    setShowWeChatLoginModal(true);
  };

  const onSubmitWeChatVerificationCode = async () => {
    const res = await API.get(
      `/api/oauth/wechat?code=${inputs.wechat_verification_code}`
    );
    const { success, message, data } = res.data;
    if (success) {
      login(data.user, data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
      localStorage.setItem('token', data.token);
      API.defaults.headers.common['Authorization'] = `Bearer ${data.token}`;
      navigate('/');
      showSuccess('登录成功！');
      setShowWeChatLoginModal(false);
    } else {
      showError(message);
    }
  };

  function handleChange(e) {
    //mock
    e.preventDefault();
    const { name, value } = e.target;
    setInputs((inputs) => ({ ...inputs, [name]: value }));
  }

  async function handleSubmit(e) {
    
    if (username && password) {
      const res = await API.post('/api/user/login', {
        username,
        password,
      });
      const { success, message, data } = res.data;
      if (success && data) {

        const { token, user } = data;

        login(user, token);
        localStorage.setItem('user', JSON.stringify(user));
        localStorage.setItem('token', token);
        // 设置 API 请求的默认 Authorization header
        API.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        // navigate('/');
        navigate('/dashboard');
        showSuccess('登录成功！');
      } else {
        showError(message);
      }
    }


    //mock login hhh
    // e.preventDefault();
    // setSubmitted(true);
    // if (username && password) {
    //   const res = await mockLogin(username, password);
    //   const { success, message, data } = res;
    //   if (success && data) {
    //     const { token, user } = data;

    //     login(user, token);
    //     localStorage.setItem('user', JSON.stringify(user));
    //     localStorage.setItem('token', token);
    //     navigate('/dashboard');
    //     showSuccess('登录成功！');
    //   } else {
    //     showError(message);
    //   }
    // }

  }

  return (
    <Grid textAlign="center" style={{ marginTop: '48px' }}>
      <Grid.Column style={{ maxWidth: 450 }}>
        <Header as="h2" color="" textAlign="center">
          <Image src="/favicon.ico" style={{ width: '40px', height: '40px' }} /> 用户登录
        </Header>
        <Form size="large">
          <Segment>
            <Form.Input
              fluid
              icon="user"
              iconPosition="left"
              placeholder="用户名"
              name="username"
              value={username}
              onChange={handleChange}
            />
            <Form.Input
              fluid
              icon="lock"
              iconPosition="left"
              placeholder="密码"
              name="password"
              type="password"
              value={password}
              onChange={handleChange}
            />
            <Button color="" fluid size="large" onClick={handleSubmit}>
              登录
            </Button>
          </Segment>
        </Form>
        <Message>
          忘记密码？
          <Link to="/reset" className="btn btn-link">
            点击重置
          </Link>
          ； 没有账户？
          <Link to="/register" className="btn btn-link">
            点击注册
          </Link>
        </Message>
        {status.github_oauth || status.wechat_login ? (
          <>
            <Divider horizontal>Or</Divider>
            {status.github_oauth ? (
              <Button
                circular
                color="black"
                icon="github"
                onClick={onGitHubOAuthClicked}
              />
            ) : (
              <></>
            )}
            {status.wechat_login ? (
              <Button
                circular
                color="green"
                icon="wechat"
                onClick={onWeChatLoginClicked}
              />
            ) : (
              <></>
            )}
          </>
        ) : (
          <></>
        )}
        <Modal
          onClose={() => setShowWeChatLoginModal(false)}
          onOpen={() => setShowWeChatLoginModal(true)}
          open={showWeChatLoginModal}
          size={'mini'}
        >
          <Modal.Content>
            <Modal.Description>
              <Image src={status.wechat_qrcode} fluid />
              <div style={{ textAlign: 'center' }}>
                <p>
                  微信扫码关注公众号，输入「验证码」获取验证码（三分钟内有效）
                </p>
              </div>
              <Form size="large">
                <Form.Input
                  fluid
                  placeholder="验证码"
                  name="wechat_verification_code"
                  value={inputs.wechat_verification_code}
                  onChange={handleChange}
                />
                <Button
                  color=""
                  fluid
                  size="large"
                  onClick={onSubmitWeChatVerificationCode}
                >
                  登录
                </Button>
              </Form>
            </Modal.Description>
          </Modal.Content>
        </Modal>
      </Grid.Column>
    </Grid>
  );
};

export default LoginForm;
