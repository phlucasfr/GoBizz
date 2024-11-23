import { useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import { maskCpfCnpj, maskPhone, validateEmail } from "../util/Index";

const SignUp = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = createSignal({
    name: "",
    email: "",
    phone: "",
    cpfCnpj: "",
    password: "",
  });
  const [error, setError] = createSignal("");

  const handleChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    let value = target.value;

    if (target.name === "phone") {
      value = maskPhone(value);
    } else if (target.name === "cpfCnpj") {
      value = maskCpfCnpj(value);
    }

    setFormData({ ...formData(), [target.name]: value });
  };

  const handleSubmit = (e: Event) => {
    e.preventDefault();
    const { name, email, phone, cpfCnpj, password } = formData();

    const cleanPhone = phone.replace(/\D/g, "");
    const cleanCpfCnpj = cpfCnpj.replace(/\D/g, "");

    if (!name || !email || !phone || !cpfCnpj || password.length < 8) {
      setError("Preencha todos os campos obrigatórios corretamente.");
      return;
    }

    if (!validateEmail(email)) {
      setError("E-mail inválido.");
      return;
    }

    setError("");
    console.log("Dados enviados:", { name, email, phone: cleanPhone, cpfCnpj: cleanCpfCnpj, password });
  };

  return (
    <div class="min-h-screen flex items-center justify-center p-4">
      <div class="w-[50vh] max-w-md p-8 rounded-lg shadow-lg space-y-6">
        <h2 class="text-3xl font-semibold text-blue-600 text-center">Cadastre Sua Empresa</h2>

        {error() && (
          <div class="mb-4 p-4 bg-red-100 text-red-800 border-l-4 border-red-500 rounded-lg">
            {error()}
          </div>
        )}

        <form onSubmit={handleSubmit} class="space-y-6">
          <div>
            <label for="name" class="block text-lg font-medium text-gray-700">Nome da Empresa</label>
            <input
              id="name"
              name="name"
              type="text"
              value={formData().name}
              onInput={handleChange}
              required
              class="mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label for="email" class="block text-lg font-medium text-gray-700">E-mail</label>
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
            <label for="phone" class="block text-lg font-medium text-gray-700">Telefone</label>
            <input
              id="phone"
              name="phone"
              type="text"
              value={formData().phone}
              onInput={handleChange}
              required
              class="mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label for="cpfCnpj" class="block text-lg font-medium text-gray-700">CPF/CNPJ</label>
            <input
              id="cpfCnpj"
              name="cpfCnpj"
              type="text"
              value={formData().cpfCnpj}
              onInput={handleChange}
              required
              class="mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label for="password" class="block text-lg font-medium text-gray-700">Senha</label>
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
            Cadastrar
          </button>
        </form>

        <div class="mt-4 text-center text-black">
          <p class="text-sm">
            Já tem uma conta?{" "}
            <a
              onClick={() => navigate("/login")}
              class="text-blue-600 hover:text-blue-500 cursor-pointer"
            >
              Faça login
            </a>
          </p>
        </div>
      </div>
    </div>
  );
};

export default SignUp;
