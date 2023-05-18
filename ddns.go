{
    order wake_on_lan  before respond
}

:2022 {
    route /re {
         respond "-"
    }
  
    route /up {
         wake_on_lan 00:11:22:33:44:55
         respond "WOL完成"
    }
  
    route / {
    	respond "{time.now}： {user_agent.name}  {user_agent.version}  {user_agent.os}  {user_agent.os_version}  {user_agent.device}   {user_agent.mobile}  {user_agent.tablet}  {user_agent.desktop}   {user_agent.bot}   {user_agent.url}"
    }
}


   
