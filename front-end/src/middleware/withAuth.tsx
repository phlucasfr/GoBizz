import { useNavigate } from "@solidjs/router";
import { validateSession } from "../api/api";
import { Component, createEffect } from "solid-js";

export const withAuth = (Component: Component) => {

  return (props: any) => {
    const navigate = useNavigate();

    createEffect(async () => {
      const sessionData = await validateSession();
      if (!sessionData.isValid) return navigate("/login", { replace: true });
    });

    return <Component {...props} />;
  };
};
