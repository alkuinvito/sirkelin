"use client";

async function getRooms() {
  return AxiosInstance.get("/api/room", { withCredentials: true });
}

export default function RoomList() {
  const { data, isLoading } = useQuery({
    queryKey: ["rooms"],
    queryFn: getRooms,
    refetchInterval: 3000,
    refetchIntervalInBackground: true,
  });

  if (isLoading) {
    return <span>Loading...</span>;
  }

  return (
    <>
      {data?.data.data.map((room) => (
        <span key={room.id}>{room.name}</span>
      ))}
    </>
  );
}
