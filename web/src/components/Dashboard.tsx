const Dashboard = (user) => {
    return (
        <div>
            <div>Dashboard {
                <ul>
                {Object.entries(user).map(([key, value]) => (
                  <li key={key}>
                    <strong>{key}:</strong> {value.toString()}
                  </li>
                ))}
              </ul>
                }
            </div>
        </div>
    )
}


export default Dashboard;