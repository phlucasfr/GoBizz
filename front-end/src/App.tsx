import { Router, Route } from "@solidjs/router";
import About from "./components/About";
import Footer from "./components/Footer";
import Header from "./components/Header";
import SignIn from "./components/SignIn";
import Welcome from "./components/Welcome";
import Contact from "./components/Contact";
import NotFound from "./components/NotFound";

const routes = [
  {
    path: "/",
    component: Welcome,
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
    path: "/contact",
    component: Contact,
  },
  {
    path: "*",
    component: NotFound,
  },
];

function App() {
  return (
    <>
      <Header />
      <Router>
        {routes.map((route) => (
          <Route path={route.path} component={route.component} />
        ))}
      </Router>
      <Footer />
    </>
  );
}

export default App;
