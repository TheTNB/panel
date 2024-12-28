import unplugin from './build/plugins/unplugin'

function genAutoImport() {
  const autoImport = unplugin[0]
  autoImport.buildStart.call({
    root: process.cwd()
  })
}

genAutoImport()
