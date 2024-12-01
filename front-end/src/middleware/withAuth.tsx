import { isLoggedIn } from "../App";
import { useNavigate } from "@solidjs/router";
import { RouteSectionProps } from "@solidjs/router";
import { JSX, createMemo, Component } from "solid-js";

export function withAuth<P extends RouteSectionProps<unknown>>(
  Component: Component<P>
) {
  return (props: P): JSX.Element => {
    const navigate = useNavigate();

    createMemo(() => {
      if (!isLoggedIn()) {
        navigate("/login", { replace: true });
      }
    });

    return <Component {...props} />;
  };
}
