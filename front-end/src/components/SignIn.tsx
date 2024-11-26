import { useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import { validateEmail } from "../util/Index";

const SignIn = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = createSignal({
    email: "",
    password: "",
  });
  const [error, setError] = createSignal("");

  const handleChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    let value = target.value;

    setFormData({ ...formData(), [target.name]: value });
  };

  const handleSubmit = (e: Event) => {
    e.preventDefault();
    const { email, password } = formData();

    if (!email || password.length < 8) {
      setError("Preencha todos os campos obrigatórios corretamente.");
      return;
    }

    if (!validateEmail(email)) {
      setError("E-mail inválido.");
      return;
    }

    setError("");
  };

  return (
    <div class="min-h-screen flex items-center justify-center p-4">
      <div class="w-[50vh] max-w-md p-8 rounded-lg shadow-lg space-y-6">
        <h2 class="text-3xl font-semibold text-blue-600 text-center">
          Você será sempre bem-vindo aqui
        </h2>

        {error() && (
          <div class="mb-4 p-4 bg-red-100 text-red-800 border-l-4 border-red-500 rounded-lg">
            {error()}
          </div>
        )}

        <form onSubmit={handleSubmit} class="space-y-6">
          <div>
            <label for="email" class="block text-lg font-medium text-gray-700">
              E-mail
            </label>
            <input
              id="email"
              name="email"
              type="email"
              value={formData().email}
              onInput={handleChange}
              required
              class="mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label
              for="password"
              class="block text-lg font-medium text-gray-700"
            >
              Senha
            </label>
            <input
              id="password"
              name="password"
              type="password"
              value={formData().password}
              onInput={handleChange}
              required
              class="mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <button
            type="submit"
            class="w-full bg-blue-600 text-white py-3 px-5 rounded-lg shadow-md hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-600"
          >
            Entrar
          </button>
        </form>

        <div class="mt-4 text-center text-black">
          <p class="text-sm">
            Não tem uma conta?{" "}
            <a
              onClick={() => navigate("/")}
              class="text-blue-600 hover:text-blue-500 cursor-pointer"
            >
              Registre-se
            </a>
          </p>
        </div>
      </div>
    </div>
  );
};

export default SignIn;
