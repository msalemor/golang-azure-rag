import axios from 'axios'
import { createEffect, createSignal, For } from 'solid-js'
import { AiOutlineCloseCircle } from 'solid-icons/ai'
import { SolidMarkdown } from 'solid-markdown'
import { Spinner, SpinnerType } from 'solid-spinner'
import { VsSend } from 'solid-icons/vs'
// import solidLogo from './assets/solid.svg'
// import viteLogo from '/vite.svg'
// import './App.css'

interface IMessage {
  role: string
  content: string
}

interface IChat {
  id: string
  messages: IMessage[]
}

interface IChats {
  chats: IChat[]
}

const Default_chat: IChat = {
  id: new Date().valueOf().toString(),
  messages: []
}

const history: IChats = {
  chats: [
    Default_chat
  ]
}

const Settings = {
  maxTokens: '',
  temperature: '0.2',
  k: '3',
  history
}

const Endpoint = import.meta.env.VITE_URL
"api/ragbot"
//const Endpoint = "http://localhost:8080/api/ragbot"

function App() {
  const [settings, setSettings] = createSignal(Settings)
  const [input, setInput] = createSignal('')
  const [process, setProcess] = createSignal(false)
  const [activeHistory, setActiveHistory] = createSignal<IChat>()

  createEffect(() => {
    setActiveHistory(settings().history.chats[settings().history.chats.length - 1])
  })

  const NewChat = () => {
    const newChat: IChat = {
      id: new Date().valueOf().toString(),
      messages: []
    }
    setSettings({ ...settings(), history: { chats: [...settings().history.chats, newChat] } })
    setActiveHistory(newChat)
  }

  const RemoveChat = (id: string) => {
    if (confirm('Are you sure you want to remove this chat?') === false) return
    setSettings({ ...settings(), history: { chats: settings().history.chats.filter(chat => chat.id !== id) } })
    setActiveHistory(settings().history.chats[settings().history.chats.length - 1])
  }

  const Process = async () => {
    if (process()) return
    if (input().length === 0) {
      alert('Input is empty')
    }
    const history: any = activeHistory()
    if (!history) {
      alert('No active chat has been selected')
      return
    }
    try {
      setProcess(true)
      const activeMessages: IMessage[] = history.messages
      activeMessages.push({ role: 'user', content: input() })
      setActiveHistory({ ...history, messages: activeMessages })

      let payload = {
        input: input(),
        messages: activeMessages,
        k: parseInt(settings().k),
        temperature: parseInt(settings().temperature),
        max_tokens: settings().maxTokens === '' ? null : parseInt(settings().maxTokens)
      }


      const resp = await axios.post(Endpoint, payload, {
        headers: {
          'Content-Type': 'application/json'
        }
      })
      const data = resp.data;
      activeMessages.push({ role: data.choices[0].message.role, content: data.choices[0].message.content })
      setActiveHistory({ ...history, messages: activeMessages })

    } catch (error) {
      console.error(error)
    } finally {
      setInput('')
      setProcess(false)
    }
  }

  return (
    <>
      <header class="text-xl font-bold bg-blue-800 dark:bg-slate-900 text-white flex items-center h-[36px] w-full px-2">
        <h1 class="">Golang RAG</h1>
      </header>
      <section class="h-[36px] bg-blue-700 dark:bg-slate-800 text-white space-x-2 flex items-center">
        <label>Settings</label>
        <label>Max Tokens:</label>
        <input class="w-20 text-black outline-none px-1" placeholder='Max'
          oninput={e => setSettings({ ...settings(), maxTokens: e.currentTarget.value })}
          value={settings().maxTokens}
        />
        <label>Temperature:</label>
        <input class="w-20 text-black outline-none px-1"
          oninput={e => setSettings({ ...settings(), temperature: e.currentTarget.value })}
          value={settings().temperature}
        />
        <label>K:</label>
        <input class="w-20 text-black outline-none px-1"
          oninput={e => setSettings({ ...settings(), k: e.currentTarget.value })}
          value={settings().k}
        />
      </section>
      <div class="flex h-[calc(100vh-36px-36px-26px)]">
        <aside class="w-[275px] bg-slate-200 dark:bg-slate-600 flex flex-col p-2 space-y-2">
          <button class='bg-blue-800 dark:bg-slate-800 p-2 w-40 text-white font-semibold hover:bg-blue-500 dark:hover:bg-slate-700'
            onclick={NewChat}
          >New Chat</button>
          <label class='bg-slate-700 text-white p-2 font-bold'>Chats</label>
          <div class='flex flex-col space-y-2'>
            <For each={settings().history.chats}>{(chat) =>
              <div class={'flex p-2 bg-blue-500 dark:bg-slate-950 text-white border-transparent hover:shadow-lg dark:hover:shadow-slate-800 hover:shadow-blue-900' + (activeHistory()?.id == chat.id ? 'border-transparent shadow-lg shadow-black dark:shadow-slate-800' : '')}
              >
                <button class='flex-grow p-1'
                  onclick={() => setActiveHistory(chat)}
                >{chat.id}</button>
                <button onClick={() => { RemoveChat(chat.id) }}><AiOutlineCloseCircle class='inline-block' /></button>
              </div>
            }</For>
          </div>
        </aside>
        <main class="w-[calc(100%-275px)] flex flex-col bg-blue-200 dark:bg-slate-700">
          <div class="flex flex-col h-[calc(100vh-36px-36px-26px-160px)] overflow-auto">
            <For each={activeHistory()?.messages}>{(message) =>
              <div class={'m-1 w-[90%] p-2 rounded ' + (message.role === 'user' ? 'ml-auto bg-blue-400 dark:bg-slate-400 text-black' : 'bg-blue-600 dark:bg-slate-600 text-white')}>
                <SolidMarkdown children={message.content} />
              </div>
            }</For>
            {process() && <Spinner type={SpinnerType.puff} class='text-blue-700 dark:text-white inline-block' height={30} />}
          </div>
          <div class="bg-blue-200 dark:bg-slate-700 flex h-[150px]">
            <textarea class="flex-grow ml-3 h-full px-1 outline-none resize-none"
              oninput={e => setInput(e.currentTarget.value)}
              value={input()}
            />
            <button class="mr-3 p-2 text-3xl hover:text-red-600 dark:hover:text-white outline-none"
              onclick={Process}
            ><VsSend class="inline-block" /></button>
          </div>
        </main>
      </div>
      <footer class="h-[26px] flex bg-blue-700 dark:bg-slate-800 text-white">
      </footer>
    </>
  )
}
export default App