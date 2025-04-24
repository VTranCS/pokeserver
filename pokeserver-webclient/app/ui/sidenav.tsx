export default function SideNav() {
    return (
        <div className="invisible md:visible bg-gray-700 w-1/2 rounded-lg">
          <div className=" h-screen relative z-0 flex h-full justify-center">
            <ul>
                <li>Guess</li>
                <li>Vote</li>
                <li>Leaderboard</li>
            </ul>
          </div>
        </div>
    );
}
