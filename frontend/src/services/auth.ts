export interface AuthUser {
  username: string;
}

interface StoredUser extends AuthUser {
  password: string;
}

const USERS_KEY = 'draft.users';
const CURRENT_USER_KEY = 'draft.currentUser';

const readUsers = (): StoredUser[] => {
  const raw = localStorage.getItem(USERS_KEY);
  if (!raw) {
    return [];
  }

  try {
    const users = JSON.parse(raw);
    return Array.isArray(users) ? users : [];
  } catch {
    return [];
  }
};

const writeUsers = (users: StoredUser[]) => {
  localStorage.setItem(USERS_KEY, JSON.stringify(users));
};

export const getCurrentUser = (): AuthUser | null => {
  const raw = localStorage.getItem(CURRENT_USER_KEY);
  if (!raw) {
    return null;
  }

  try {
    const user = JSON.parse(raw);
    if (!user?.username) {
      return null;
    }
    return { username: user.username };
  } catch {
    return null;
  }
};

export const register = (username: string, password: string): AuthUser => {
  const normalizedUsername = username.trim();
  if (!normalizedUsername || !password) {
    throw new Error('用户名和密码不能为空');
  }

  const users = readUsers();
  const exists = users.some((u) => u.username === normalizedUsername);
  if (exists) {
    throw new Error('用户已存在');
  }

  users.push({ username: normalizedUsername, password });
  writeUsers(users);

  return { username: normalizedUsername };
};

export const login = (username: string, password: string): AuthUser => {
  const normalizedUsername = username.trim();
  if (!normalizedUsername || !password) {
    throw new Error('用户名和密码不能为空');
  }

  const users = readUsers();
  const user = users.find(
    (u) => u.username === normalizedUsername && u.password === password,
  );

  if (!user) {
    throw new Error('用户名或密码错误');
  }

  const currentUser: AuthUser = { username: user.username };
  localStorage.setItem(CURRENT_USER_KEY, JSON.stringify(currentUser));
  return currentUser;
};

export const logout = () => {
  localStorage.removeItem(CURRENT_USER_KEY);
};

