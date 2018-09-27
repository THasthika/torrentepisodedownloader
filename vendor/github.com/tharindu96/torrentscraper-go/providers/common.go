package providers

// Provider Implementation
type Provider struct {
	ID             string
	Name           string
	SupportedTypes TorrentType
	Search         func(query string, t TorrentType, out chan []*TorrentMeta)
	SearchShow     func(name string, season uint, episode uint, out chan []*TorrentMeta)
	SearchMovie    func(name string, out chan []*TorrentMeta)
}

// TorrentType of the Torrent
type TorrentType uint

// Set of TorrentTypes
const (
	TorrentTypeUnspecified = 0x00
	TorrentTypeMovie       = 0x01
	TorrentTypeTV          = 0x02
	TorrentTypeGame        = 0x04
	TorrentTypeBook        = 0x08
)

// TorrentMeta stores the information about the torrent
type TorrentMeta struct {
	Name   string
	Magnet string
	Seeds  uint
	Type   TorrentType
	Size   uint
}
