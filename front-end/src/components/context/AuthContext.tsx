import { createContext, useContext, createSignal, JSX } from "solid-js";

type AuthContextType = {
  login: () => void;
  logout: () => void;
  isLoggedIn: () => boolean;
};

const AuthContext = createContext<AuthContextType>();

export const AuthProvider = (props: {
  children:
    | number
    | boolean
    | Node
    | JSX.ArrayElement
    | (string & {})
    | null
    | undefined;
}) => {
  
  const [isLoggedIn, setIsLoggedIn] = createSignal(false);

  const login = () => setIsLoggedIn(true);
  const logout = () => setIsLoggedIn(false);

  return (
    <AuthContext.Provider value={{ isLoggedIn, login, logout }}>
      {props.children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

