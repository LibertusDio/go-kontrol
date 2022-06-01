import React from "react";

const Home: React.FC = () => {
  return (
      <div>
          <div className="buttons">
              <a href={'/idt/service'} className="button is-info">idt</a>
              <a href={'/adt/service'} className="button is-success">adt</a>
              <a href={'/hrd/service'} className="button is-warning">hrd</a>
          </div>
      </div>
  )
}

export default Home
