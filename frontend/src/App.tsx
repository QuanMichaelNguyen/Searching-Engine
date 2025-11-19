import "./App.css";

function App() {
  return (
    <>
      {/* --- The Title --- */}
      <div className="flex flex-col gap-40 items-center mt-28">
        <h1 className="text-7xl font-light tracking-wide text-transparent bg-clip-text bg-gradient-to-r from-blue-600 via-purple-500 to-pink-500 transform scale-y-100 hover:scale-y-105 transition-transform duration-500 ease-in-out cursor-default">
          GoFind
        </h1>

        {/* --- The Search Bar --- */}
        <div className="relative w-full max-w-xl">
          <svg
            className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            ></path>
          </svg>
          <input
            type="text"
            placeholder="Enter your search"
            className="search-input w-full py-3 pl-12 pr-4 text-lg border-2 border-gray-300 rounded-full shadow-lg focus:outline-none focus:ring-4 focus:ring-blue-200 focus:border-blue-500 transition duration-300 ease-in-out"
          />
        </div>
      </div>
    </>
  );
}

export default App;
