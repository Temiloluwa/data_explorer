export const HomeFooter = () => {
  return (
    <div className="relative flex h-8 w-full">
      <div className="absolute bottom-0 left-0 right-0 bg-slate-500">
        <div className=" border-t p-4 flex items-center gap-2 max-w-4xl mx-auto w-full">
          <input
            type="text"
            placeholder="Type a message..."
            className="flex-1 border p-2 rounded"
          />
          <button className="bg-blue-500 text-white px-4 py-2 rounded">
            Send
          </button>
        </div>
      </div>
    </div>
  );
};
