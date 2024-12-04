import { Motion } from "@motionone/solid";
import { useAuth } from "./context/AuthContext";
import { Home, Info, Mail, Menu, X } from "lucide-solid";
import { createSignal, createMemo, For, Show } from "solid-js";

type SectionKey = "about" | "contact" | "login" | "register";
type Section = {
  id: SectionKey;
  href: string;
  icon: typeof Home;
  title: string;
};

const Header = () => {
  const { isLoggedIn } = useAuth();
  const [isMenuOpen, setIsMenuOpen] = createSignal(false);

  const closeMenu = () => setIsMenuOpen(false);

  const sections = createMemo<Section[]>(() => [
    {
      id: "about",
      href: "/about",
      icon: Info,
      title: "Sobre Nós",
    },
    {
      id: "contact",
      href: "/contact",
      icon: Mail,
      title: "Contato",
    },
  ]);

  return (
    <Show when={!isLoggedIn()}>
      <Motion.div
        initial={{ opacity: 0, y: -50 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        class="fixed top-0 left-0 right-0 z-50 bg-white shadow-md"
      >
        <nav class="container mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex items-center justify-between h-16">
            <div class="flex items-center">
              <a href="/" class="flex-shrink-0">
                <img
                  class="h-10 w-auto"
                  src="/assets/logo.png"
                  alt="GoBizz Logo"
                />
              </a>
              <div class="hidden md:block ml-10">
                <div class="flex items-baseline space-x-4">
                  <For each={sections()}>
                    {(section) => (
                      <a
                        href={section.href}
                        class="text-gray-600 hover:text-indigo-600 px-3 py-2 rounded-md text-sm font-medium transition-colors duration-300"
                      >
                        {section.title}
                      </a>
                    )}
                  </For>
                </div>
              </div>
            </div>
            <div class="hidden md:flex items-center space-x-4">
              <a
                href="/register"
                class="bg-white text-indigo-600 border border-indigo-600 px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-50 transition-colors duration-300"
              >
                Cadastrar
              </a>
              <a
                href="/login"
                class="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700 transition-colors duration-300"
              >
                Entrar
              </a>
            </div>
            <div class="-mr-2 flex md:hidden">
              <Motion.button
                onClick={() => setIsMenuOpen(!isMenuOpen())}
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-indigo-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500"
                aria-expanded="false"
              >
                <span class="sr-only">Abrir menu principal</span>
                {isMenuOpen() ? (
                  <X class="block h-6 w-6" aria-hidden="true" />
                ) : (
                  <Menu class="block h-6 w-6" aria-hidden="true" />
                )}
              </Motion.button>
            </div>
          </div>
        </nav>

        <Motion.div
          initial={{ opacity: 0, height: 0 }}
          animate={{
            opacity: isMenuOpen() ? 1 : 0,
            height: isMenuOpen() ? "auto" : 0,
          }}
          transition={{ duration: 0.3 }}
          class={`md:hidden ${isMenuOpen() ? "block" : "hidden"}`}
        >
          <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
            <For each={sections()}>
              {(section) => (
                <a
                  href={section.href}
                  class="text-gray-600 hover:text-indigo-600 block px-3 py-2 rounded-md text-base font-medium transition-colors duration-300"
                  onClick={closeMenu}
                >
                  {section.title}
                </a>
              )}
            </For>
            <a
              href="/register"
              class="block w-full text-left bg-white text-indigo-600 border border-indigo-600 px-3 py-2 rounded-md text-base font-medium hover:bg-indigo-50 transition-colors duration-300 mt-2"
              onClick={closeMenu}
            >
              Cadastrar
            </a>
            <a
              href="/login"
              class="block w-full text-left bg-indigo-600 text-white px-3 py-2 rounded-md text-base font-medium hover:bg-indigo-700 transition-colors duration-300 mt-2"
              onClick={closeMenu}
            >
              Entrar
            </a>
          </div>
        </Motion.div>
      </Motion.div>
      <div class="h-16"></div>
    </Show>
  );
};

export default Header;

