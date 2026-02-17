import React, { useState, useEffect } from 'react';
import Dashboard from './pages/Dashboard';
import DictionaryUpload from './pages/DictionaryUpload';
import Learning from './pages/Learning';
import Hello from './pages/Hello';
import {
  clearAuth,
  fetchCurrentUser,
  getCurrentUser,
  logout,
  type AuthUser,
} from './services/auth';
import './styles/global.css';

type Page = 'hello' | 'dashboard' | 'upload' | 'learn';

const App: React.FC = () => {
  const [currentPage, setCurrentPage] = useState<Page>('hello');
  const [dictId, setDictId] = useState<number>(1);
  const [user, setUser] = useState<AuthUser | null>(null);

  const getPageFromPath = (path: string): Page => {
    if (path === '/hello') {
      return 'hello';
    }
    if (path === '/upload') {
      return 'upload';
    }
    if (path === '/learn') {
      return 'learn';
    }
    return 'dashboard';
  };

  const guardPage = (targetPage: Page, currentUser: AuthUser | null): Page => {
    if (!currentUser && targetPage !== 'hello') {
      return 'hello';
    }
    if (currentUser && targetPage === 'hello') {
      return 'dashboard';
    }
    return targetPage;
  };

  useEffect(() => {
    const syncRouteFromLocation = () => {
      const cachedUser = getCurrentUser();
      setUser(cachedUser);

      const path = window.location.pathname;
      const searchParams = new URLSearchParams(window.location.search);
      const targetPage = getPageFromPath(path);
      const finalPage = guardPage(targetPage, cachedUser);

      if (finalPage === 'learn') {
        const id = searchParams.get('dictId');
        if (id) {
          setDictId(parseInt(id, 10));
        }
      }

      setCurrentPage(finalPage);

      if (finalPage !== targetPage) {
        const redirectPath = finalPage === 'hello' ? '/hello' : '/';
        window.history.replaceState({}, '', redirectPath);
      }
    };

    syncRouteFromLocation();
    fetchCurrentUser()
      .then((freshUser) => {
        setUser(freshUser);
      })
      .catch(() => {
        clearAuth();
        setUser(null);
        setCurrentPage('hello');
        window.history.replaceState({}, '', '/hello');
      });
    window.addEventListener('popstate', syncRouteFromLocation);
    return () => window.removeEventListener('popstate', syncRouteFromLocation);
  }, []);

  const navigateTo = (
    page: Page,
    params?: { dictId?: number },
    replaceHistory = false,
    userOverride?: AuthUser | null,
  ) => {
    const activeUser = userOverride === undefined ? user : userOverride;
    const finalPage = guardPage(page, activeUser);

    setCurrentPage(finalPage);
    if (params?.dictId) {
      setDictId(params.dictId);
    }

    // Update URL
    let url = '/hello';
    if (finalPage === 'dashboard') {
      url = '/';
    } else if (finalPage === 'upload') {
      url = '/upload';
    } else if (finalPage === 'learn') {
      url = params?.dictId ? `/learn?dictId=${params.dictId}` : '/learn';
    }

    if (replaceHistory) {
      window.history.replaceState({}, '', url);
    } else {
      window.history.pushState({}, '', url);
    }
  };

  // Override window.location for navigation
  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      const anchor = target.closest('a');
      
      if (anchor) {
        const href = anchor.getAttribute('href');
        if (href && href.startsWith('/')) {
          e.preventDefault();
          if (href === '/hello') {
            navigateTo('hello');
          } else if (href === '/upload') {
            navigateTo('upload');
          } else if (href.startsWith('/learn')) {
            const url = new URL(href, window.location.origin);
            const dictId = url.searchParams.get('dictId');
            navigateTo('learn', { dictId: dictId ? parseInt(dictId, 10) : 1 });
          } else {
            navigateTo('dashboard');
          }
        }
      }
    };

    document.addEventListener('click', handleClick);
    return () => document.removeEventListener('click', handleClick);
  }, [user]);

  const handleLoginSuccess = (loggedInUser: AuthUser) => {
    setUser(loggedInUser);
    navigateTo('dashboard', undefined, true, loggedInUser);
  };

  const handleLogout = async () => {
    await logout();
    setUser(null);
    navigateTo('hello', undefined, true, null);
  };

  const renderPage = () => {
    switch (currentPage) {
      case 'hello':
        return <Hello onLoginSuccess={handleLoginSuccess} />;
      case 'upload':
        return <DictionaryUpload />;
      case 'learn':
        return <Learning dictId={dictId} />;
      case 'dashboard':
      default:
        return <Dashboard />;
    }
  };

  return (
    <div className="app">
      {user ? (
        <button
          type="button"
          onClick={handleLogout}
          style={{
            position: 'fixed',
            right: '16px',
            top: '16px',
            zIndex: 100,
            border: 'none',
            borderRadius: '10px',
            padding: '8px 12px',
            background: '#1f2937',
            color: '#fff',
            cursor: 'pointer',
          }}
        >
          退出登录（{user.username}）
        </button>
      ) : null}
      {renderPage()}
    </div>
  );
};

export default App;
