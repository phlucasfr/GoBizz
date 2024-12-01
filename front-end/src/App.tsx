import Home from "./components/Home";
import About from "./components/About";
import Footer from "./components/Footer";
import Header from "./components/Header";
import SignIn from "./components/SignIn";
import Welcome from "./components/Welcome";
import Contact from "./components/Contact";
import NotFound from "./components/NotFound";
import { withAuth } from "./middleware/withAuth";
import { Router, Route } from "@solidjs/router";
import { validateSession } from "./api/api";
import { createSignal, onMount, Show } from "solid-js";

export const [isLoggedIn, setIsLoggedIn] = createSignal(false);
export const [loadingAuth, setLoadingAuth] = createSignal(true);

const App = () => {
  onMount(async () => {
    const sessionData = await validateSession().then();

    if (sessionData instanceof Error) {
      if (sessionData.message !== "Usuário não autorizado") {
        console.error(sessionData);
      }
      setIsLoggedIn(false);
    } else {
      setIsLoggedIn(sessionData.isValid);
    }

    setLoadingAuth(false);
  });

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
      path: "/contact",
      component: Contact,
    },
    {
      path: "*",
      component: NotFound,
    },
  ];

  return (
    <Show
      when={!loadingAuth()}
      fallback={<div>Carregando autenticação...</div>}
    >
      <>
        <Header />
        <Router>
          {routes.map((route) => (
            <Route path={route.path} component={route.component} />
          ))}
        </Router>
        <Footer />
      </>
    </Show>
  );
};

export default App;
