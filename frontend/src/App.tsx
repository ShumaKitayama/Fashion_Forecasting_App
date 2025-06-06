import { Routes, Route } from "react-router-dom";

function App() {
  return (
    <div className="app-container">
      <header className="app-header">
        <h1>TrendScout</h1>
      </header>
      <main>
        <Routes>
          <Route
            path="/"
            element={<div>Welcome to TrendScout - Dashboard coming soon</div>}
          />
          <Route path="/login" element={<div>Login Page - Coming soon</div>} />
          <Route
            path="/register"
            element={<div>Register Page - Coming soon</div>}
          />
          <Route path="*" element={<div>404 - Page not found</div>} />
        </Routes>
      </main>
      <footer className="app-footer">
        <p>
          &copy; {new Date().getFullYear()} TrendScout - All rights reserved
        </p>
      </footer>
    </div>
  );
}

export default App;
