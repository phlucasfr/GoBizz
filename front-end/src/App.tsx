import Home from "./components/Home";
import About from "./components/About";
import Footer from "./components/Footer";
import Header from "./components/Header";
import SignIn from "./components/SignIn";
import SignUp from "./components/SignUp";
import Welcome from "./components/Welcome";
import Contact from "./components/Contact";
import NotFound from "./components/NotFound";
import ResetPasswordPage from "./components/pages/ResetPasswordPage";
import EmailVerificationPage from "./components/pages/EmailVerificationPage";

import { withAuth } from "./middleware/withAuth";
import { AuthProvider } from "./components/context/AuthContext";
import { Router, Route } from "@solidjs/router";

const App = () => {
  const routes = [
    {
      path: "/",
      component: Welcome,
    },
    {
      path: "/home/:id",
      component: withAuth(Home),
    },   
    {
      path: "/about",
      component: About,
    },
    {
      path: "/login",
      component: SignIn,
    },
    {
      path: "/register",
      component: SignUp,
    },
    {
      path: "/contact",
      component: Contact,
    },
    {
      path: "*",
      component: NotFound,
    },
    {
      path: "/reset-password",
      component: ResetPasswordPage
    },
    {
      path: "/email-verification",
      component: EmailVerificationPage
    },
  ];

  return (
    <>
      <AuthProvider>
        <Header />
        <Router>
          {routes.map((route) => (
            <Route path={route.path} component={route.component} />
          ))}
        </Router>
        <Footer />
      </AuthProvider>
    </>
  );
};

export default App;
