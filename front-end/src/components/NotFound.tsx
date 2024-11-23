import { useNavigate } from "@solidjs/router";

const NotFound = () => {
  const navigate = useNavigate();

  return (
    <div class="min-h-[60vh] flex flex-col items-center justify-center text-center px-4">
      <h1 class="text-6xl font-bold text-gray-900 mb-4">404</h1>
      <p class="text-xl text-gray-600 mb-8">Oops! Página não encontrada.</p>
      <p class="text-gray-500 mb-8">
        A página que você está procurando não existe ou foi movida.
      </p>
      <button
        onClick={() => navigate("/")}
        class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-300"
      >
        Voltar para Home
      </button>
    </div>
  );
};

export default NotFound;