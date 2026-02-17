import React, { useState } from 'react';
import { login, register, type AuthUser } from '../../services/auth';
import './Hello.css';

interface HelloPageProps {
  onLoginSuccess: (user: AuthUser) => void;
}

const Hello: React.FC<HelloPageProps> = ({ onLoginSuccess }) => {
  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const resetError = () => {
    if (error) {
      setError('');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      let user: AuthUser;
      if (mode === 'register') {
        user = await register(username, password);
      } else {
        user = await login(username, password);
      }
      onLoginSuccess(user);
    } catch (err) {
      const message = err instanceof Error ? err.message : '操作失败，请重试';
      setError(message);
    }
  };

  return (
    <div className="hello-page">
      <div className="hello-card">
        <h1>hello</h1>
        <p className="hello-subtitle">请先登录后再访问应用功能</p>

        <div className="hello-tabs">
          <button
            className={mode === 'login' ? 'active' : ''}
            onClick={() => {
              setMode('login');
              resetError();
            }}
            type="button"
          >
            登录
          </button>
          <button
            className={mode === 'register' ? 'active' : ''}
            onClick={() => {
              setMode('register');
              resetError();
            }}
            type="button"
          >
            注册
          </button>
        </div>

        <form onSubmit={handleSubmit} className="hello-form">
          <label>
            用户名
            <input
              type="text"
              value={username}
              onChange={(e) => {
                setUsername(e.target.value);
                resetError();
              }}
              placeholder="请输入用户名"
            />
          </label>
          <label>
            密码
            <input
              type="password"
              value={password}
              onChange={(e) => {
                setPassword(e.target.value);
                resetError();
              }}
              placeholder="请输入密码"
            />
          </label>

          {error ? <p className="hello-error">{error}</p> : null}

          <button className="hello-submit" type="submit">
            {mode === 'login' ? '登录' : '注册并登录'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default Hello;
