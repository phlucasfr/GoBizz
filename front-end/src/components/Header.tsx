import type { JSX } from "solid-js";
import { createSignal } from "solid-js";

type SectionKey = "home" | "about" | "contact";
type Section = {
  title: string;
  icon: () => JSX.Element;
  href: string;
};

const Header = () => {
  const sections: Record<SectionKey, Section> = {
    home: {
      title: "Home",
      icon: () => (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          class="w-5 h-5"
        >
          <path
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 2L2 7h3v8h2V9h6v6h2V7h3L12 2z"
          />
        </svg>
      ),
      href: "/",
    },
    about: {
      title: "Sobre Nós",
      icon: () => (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          class="w-5 h-5"
        >
          <path
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 4v7l3-3h3a1 1 0 011 1v9a1 1 0 01-1 1H5a1 1 0 01-1-1V9a1 1 0 011-1h3l3 3V4"
          />
        </svg>
      ),
      href: "/about",
    },
    contact: {
      title: "Contato",
      icon: () => (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          class="w-5 h-5"
        >
          <path
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M22 2v20H2V2h20zm-2 2H4v16h16V4z"
          />
        </svg>
      ),
      href: "/contact",
    },
  };

  const [isMenuOpen, setIsMenuOpen] = createSignal(false);

  return (
    <div class="min-h-20">
      <header class="fixed top-0 left-0 right-0 z-50 shadow-lg">
        <div class="bg-gradient-to-r from-blue-600/70 to-purple-600/70 text-white rounded-b-md">
          <div class="container mx-auto px-6 py-4 flex items-center justify-between">
            <button
              class="md:hidden focus:outline-none rounded-full p-2 hover:bg-blue/20"
              onClick={() => setIsMenuOpen(!isMenuOpen())}
              aria-label="Toggle Menu"
            >
              {isMenuOpen() ? (
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  class="w-6 h-6"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              ) : (
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  class="w-6 h-6"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 6h16M4 12h16M4 18h16"
                  />
                </svg>
              )}
            </button>

            <nav
              class={`
                ${isMenuOpen() ? "block" : "hidden"} 
                absolute top-full left-0 w-full bg-gradient-to-r from-indigo-400/60 to-purple-600/60
                md:static md:block md:w-auto
                transition-all duration-300 ease-in-out
                rounded-3xl shadow-lg
              `}
            >
              <ul class="flex flex-col md:flex-row space-y-4 md:space-y-0 md:space-x-6">
                {Object.keys(sections).map((key) => (
                  <li
                    class={`
                      group relative cursor-pointer 
                      hover:bg-purple-700/20
                      rounded-full transition-transform duration-300 ease-in-out
                      hover:scale-110
                    `}
                  >
                    <a
                      class="
                        w-full flex items-center justify-center md:justify-start 
                        gap-3 px-5 py-3 text-sm font-medium rounded-full
                        bg-transparent hover:bg-white/20 transition-colors
                      "
                      href={sections[key as SectionKey].href}
                    >
                      {sections[key as SectionKey].icon()}
                      <span class="hidden md:inline">
                        {sections[key as SectionKey].title}
                      </span>
                    </a>
                  </li>
                ))}
              </ul>
            </nav>
          </div>
        </div>
      </header>
    </div>
  );
};

export default Header;
