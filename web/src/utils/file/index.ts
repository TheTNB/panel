const getExt = (filename: string) => {
  const dot = filename.lastIndexOf('.')
  if (dot === -1 || dot === 0) {
    return ''
  }
  return filename.slice(dot + 1)
}

const getBase = (filename: string) => {
  const dot = filename.lastIndexOf('.')
  if (dot === -1 || dot === 0) {
    return filename
  }
  return filename.slice(0, dot)
}

const getIconByExt = (ext: string) => {
  switch (ext) {
    case 'png':
    case 'jpg':
    case 'jpeg':
    case 'gif':
      return 'bi:file-earmark-image'
    case 'mp4':
    case 'avi':
    case 'mkv':
    case 'rmvb':
      return 'bi:file-earmark-play'
    case 'mp3':
    case 'flac':
    case 'wav':
    case 'ape':
      return 'bi:file-earmark-music'
    case 'zip':
    case 'rar':
    case '7z':
    case 'tar':
    case 'gz':
      return 'bi:file-earmark-zip'
    case 'doc':
    case 'docx':
    case 'xls':
    case 'xlsx':
      return 'bi:file-earmark-word'
    case 'ppt':
    case 'pptx':
      return 'bi:file-earmark-ppt'
    case 'pdf':
      return 'bi:file-earmark-pdf'
    case 'txt':
    case 'md':
    case 'log':
    case 'conf':
    case 'ini':
    case 'yaml':
    case 'yml':
      return 'bi:file-earmark-text'
    case 'html':
    case 'htm':
    case 'xml':
    case 'json':
    case 'js':
    case 'css':
    case 'ts':
    case 'vue':
    case 'jsx':
    case 'tsx':
    case 'php':
    case 'java':
    case 'py':
    case 'go':
    case 'rb':
    case 'sh':
      return 'bi:file-earmark-code'
    case '':
      return 'bi:file-earmark-binary'
    default:
      return 'bi:file-earmark'
  }
}

const languageByPath = (path: string) => {
  if (path.startsWith('/www/server/openresty/')) {
    return 'nginx'
  }

  const ext = getExt(path)
  switch (ext) {
    case 'abap':
      return 'abap'
    case 'apex':
      return 'apex'
    case 'azcli':
      return 'azcli'
    case 'bat':
      return 'bat'
    case 'bicep':
      return 'bicep'
    case 'mligo': // cameligo 扩展名
      return 'cameligo'
    case 'clj':
    case 'cljs':
    case 'cljc': // clojure 扩展名
      return 'clojure'
    case 'coffee':
      return 'coffee'
    case 'cpp':
    case 'cc':
    case 'cxx': // cpp 扩展名
      return 'cpp'
    case 'cs':
      return 'csharp'
    case 'csp':
      return 'csp'
    case 'css':
      return 'css'
    case 'cypher':
      return 'cypher'
    case 'dart':
      return 'dart'
    case 'dockerfile':
      return 'dockerfile'
    case 'ecl':
      return 'ecl'
    case 'ex':
    case 'exs': // elixir 扩展名
      return 'elixir'
    case 'flow':
      return 'flow9'
    case 'fs':
    case 'fsi':
    case 'fsx':
    case 'fsscript': // fsharp 扩展名
      return 'fsharp'
    case 'ftl': // freemarker2 扩展名
      return 'freemarker2'
    case 'go':
      return 'go'
    case 'graphql':
      return 'graphql'
    case 'handlebars':
    case 'hbs': // handlebars 扩展名
      return 'handlebars'
    case 'hcl':
    case 'tf': // hcl 扩展名
      return 'hcl'
    case 'html':
    case 'htm': // html 扩展名
      return 'html'
    case 'ini':
      return 'ini'
    case 'java':
      return 'java'
    case 'js':
    case 'mjs':
    case 'cjs': // javascript 扩展名
      return 'javascript'
    case 'jl':
      return 'julia'
    case 'kt':
    case 'kts': // kotlin 扩展名
      return 'kotlin'
    case 'less':
      return 'less'
    case 'lex': // lexon 扩展名
      return 'lexon'
    case 'lua':
      return 'lua'
    case 'liquid':
      return 'liquid'
    case 'm3':
      return 'm3'
    case 'md':
      return 'markdown'
    case 'mdx':
      return 'mdx'
    case 'mips':
      return 'mips'
    case 'dax': // msdax 扩展名
      return 'msdax'
    case 'm': // objective-c 扩展名
      return 'objective-c'
    case 'pas':
      return 'pascal'
    case 'ligo': // pascaligo 扩展名
      return 'pascaligo'
    case 'pl':
      return 'perl'
    case 'php':
      return 'php'
    case 'pla':
      return 'pla'
    case 'dats':
    case 'sats':
    case 'hats': // postiats 扩展名
      return 'postiats'
    case 'pq': // powerquery 扩展名
      return 'powerquery'
    case 'ps1':
    case 'psm1': // powershell 扩展名
      return 'powershell'
    case 'proto':
      return 'protobuf'
    case 'pug':
    case 'jade': // pug 扩展名
      return 'pug'
    case 'py':
      return 'python'
    case 'qs': // qsharp 扩展名
      return 'qsharp'
    case 'r':
      return 'r'
    case 'razor':
      return 'razor'
    case 'redis':
      return 'redis'
    case 'redshift':
      return 'redshift'
    case 'rst':
      return 'restructuredtext'
    case 'rb':
      return 'ruby'
    case 'rs':
      return 'rust'
    case 'sb':
      return 'sb'
    case 'scala':
      return 'scala'
    case 'scm':
    case 'ss': // scheme 扩展名
      return 'scheme'
    case 'scss':
      return 'scss'
    case 'sh':
    case 'bash': // shell 扩展名
      return 'shell'
    case 'sol':
      return 'solidity'
    case 'sophia':
      return 'sophia'
    case 'sparql':
      return 'sparql'
    case 'sql':
      return 'sql'
    case 'st':
      return 'st'
    case 'swift':
      return 'swift'
    case 'sv':
    case 'svh': // systemverilog 扩展名
      return 'systemverilog'
    case 'tcl':
      return 'tcl'
    case 'twig':
      return 'twig'
    case 'ts':
    case 'tsx': // typescript 扩展名
      return 'typescript'
    case 'typespec':
      return 'typespec'
    case 'vb':
    case 'vbs': // vb 扩展名
      return 'vb'
    case 'wgsl':
      return 'wgsl'
    case 'xml':
    case 'xsd':
    case 'xsl':
    case 'xslt': // xml 扩展名
      return 'xml'
    case 'yaml':
    case 'yml': // yaml 扩展名
      return 'yaml'
    default:
      return ''
  }
}

const checkName = (name: string) => {
  return /^[a-zA-Z0-9_.@#$%\-\s[\]()]+$/.test(name)
}

const checkPath = (path: string) => {
  return /^(?!\/)(?!.*\/$)(?!.*\/\/)(?!.*\s).*$/.test(path)
}

const getFilename = (path: string) => {
  const parts = path.split('/')
  return parts.pop()!
}

const isArchive = (name: string) => {
  const ext = getExt(name)
  return ['zip', 'rar', '7z', 'tar', 'gz'].includes(ext)
}

const formatPercent = (num: any) => {
  num = Number(num)
  return Number(num.toFixed(2))
}

const formatBytes = (size: any) => {
  size = Number(size)
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  let i = 0

  while (size >= 1024 && i < units.length) {
    size /= 1024
    i++
  }

  return size.toFixed(2) + ' ' + units[i]
}

const lastDirectory = (path: string) => {
  const parts = path.split('/')
  return parts.pop() || ''
}

export {
  getExt,
  getBase,
  getIconByExt,
  languageByPath,
  checkName,
  checkPath,
  getFilename,
  isArchive,
  formatPercent,
  formatBytes,
  lastDirectory
}
