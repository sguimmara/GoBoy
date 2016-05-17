package lcd

import (
	. "common"
	"cpu"
	"github.com/veandco/go-sdl2/sdl"
	// "log"
	mem "memory"
	"os"
	"runtime/pprof"
	"time"
	"unsafe"
)

const (
	WHITE      uint32 = 0xFFFFFFFF
	DARK_GREY  uint32 = 0x44444444
	LIGHT_GREY uint32 = 0xAAAAAAAA
	BLACK      uint32 = 0x00000000
	LCD_WIDTH         = 160
	LCD_HEIGHT        = 144
	SCANLINES         = 153

	SCY_ADDR     = 0xFF42
	SCX_ADDR     = 0xFF43
	OAM_ADDR     = 0xFE00
	OAM_ADDR_END = 0xFE9F
	WY_ADDR      = 0xFF4A
	WX_ADDR      = 0xFF4B
	STAT_ADDR    = 0xFF41

	STAT_LYC     = 6
	STAT_MODE_10 = 5
	STAT_MODE_01 = 4
	STAT_MODE_00 = 3
	STAT_COIN    = 2

	LCD_ACTIVE    = 7
	WDW_MAP       = 6
	WDW_ACTIVE    = 5
	TDT           = 4
	BG_MAP        = 3
	SPRITE_SIZE   = 2
	SPRITE_ACTIVE = 1
	BG_WDW_ACTIVE = 0

	MODE_0_HBLANK        = 0
	MODE_1_VBLANK        = 1
	MODE_2_OAM_USED      = 2
	MODE_3_OAM_VRAM_USED = 3

	SCALE = 3
)

var (
	window    *sdl.Window
	renderer  *sdl.Renderer
	screenTex *sdl.Texture
	pixels    = new([LCD_HEIGHT * LCD_WIDTH]uint32)

	palette = new([4]uint32)

	mmap     []byte = nil
	CONTINUE        = true
)

func init() {
	palette[0] = WHITE
	palette[1] = LIGHT_GREY
	palette[2] = DARK_GREY
	palette[3] = BLACK

	// To reduce overhead, lcd code can directly manipulate RAM.
	mmap = mem.GetMemoryMap()
}

func Initialize() {
	var err error

	sdl.Init(sdl.INIT_EVERYTHING)

	window, err = sdl.CreateWindow("goboy", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		LCD_WIDTH*SCALE, LCD_HEIGHT*SCALE, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	screenTex, err = renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, LCD_WIDTH, LCD_HEIGHT)
	if err != nil {
		panic(err)
	}

	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Clear()
}

func RunProfile() {
	f, _ := os.Create("cpuprof.txt")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < 1000; i++ {
		redraw()
	}
}

func Run() {
	CONTINUE = true
	for {
		if CONTINUE {
			redraw()
		}
	}
}

func Pause() {
	CONTINUE = false
}

func Continue() {
	CONTINUE = true
}

func Stop() {
	screenTex.Destroy()
	window.Destroy()
	sdl.Quit()
}

func SetTile(index int, tile []byte, base uint16) {
	offset := uint16(len(tile)) * uint16(index)
	mem.SetRange(base+offset, tile)
}

func SetBackgroundTile(x, y, index int) {
	addr := uint16(0x9800)
	if IsBitSet(mem.GetLCDC(), BG_MAP) {
		addr = uint16(0x9C00)
	}

	mem.Set(addr+uint16(y*32+x), uint8(index))
}

func SetWindowTile(x, y, index int) {
	addr := uint16(0x9800)
	if IsBitSet(mem.GetLCDC(), WDW_MAP) {
		addr = uint16(0x9C00)
	}

	mem.Set(addr+uint16(y*32+x), uint8(index))
}

// Create a sprite at index x (FFE0 + x), of coordinates x, y,
// using pattern <pattern> and with flags <flags>
func SetSprite(index uint16, x, y, pattern, flags uint8) {
	if index > 0x9F {
		panic("trying to set a sprite outside of OAM range.")
	}

	base := OAM_ADDR + index*4
	mmap[base] = flags
	mmap[base+1] = pattern
	mmap[base+2] = y
	mmap[base+3] = x
}

func setBufferPixel(x, y, color int, pixels unsafe.Pointer, pitch int) {
	(*[LCD_WIDTH * LCD_HEIGHT]uint32)(pixels)[y*(pitch/4)+x] = palette[color]
}

func getBackgroundTileMap() uint16 {
	addr := uint16(0x9800)
	if IsBitSet(mem.GetLCDC(), BG_MAP) {
		addr = uint16(0x9C00)
	}

	return addr
}

func getWindowTileMap() uint16 {
	addr := uint16(0x9800)
	if IsBitSet(mem.GetLCDC(), WDW_MAP) {
		addr = uint16(0x9C00)
	}

	return addr
}

func getTileDataTable() uint16 {
	tdt := uint16(0x8800)
	if IsBitSet(mem.GetLCDC(), TDT) {
		tdt = uint16(0x8000)
	}
	return tdt
}

func getSpriteColor(x, y int) int {
	var col int
	for s := OAM_ADDR; s <= OAM_ADDR_END; s += 4 {
		// get the sprite coordinates
		sx := int(mmap[s+3])
		sy := int(mmap[s+2])

		pattern := int(mmap[s+1])
		flags := uint8(mmap[s])

		flipX := IsBitSet(flags, 5)
		flipY := IsBitSet(flags, 6)
		fx := x - sx + 8
		fy := y - sy + 8

		if flipX {
			fx = 7 - fx
		}
		if flipY {
			fy = 7 - fy
		}

		if x >= sx-8 && x < sx && y >= sy-8 && y < sy {
			col = getSpritePixel(fx, fy, pattern)
			if col != 0 {
				// break early so that the first sprite found has priority
				// Normally, priority determination is a bit more complex, but
				// for now, it will do.
				break
			}
		}
	}
	return col
}

func drawWindowLine(y int, mapAddr, tileAddr uint16, pixels unsafe.Pointer, pitch int) {
	x := int(mem.GetWX())
	yy := y + int(mem.GetWY())

	for i := 0; i < LCD_WIDTH; i++ {
		pix := getTileColor(x, yy, mapAddr, tileAddr)
		setBufferPixel(x, y, pix, pixels, pitch)
		x++
	}
}

// Draw a single line of background
func drawBackgroundLine(y int, mapAddr, tileAddr uint16, pixels unsafe.Pointer, pitch int) {
	x := int(mmap[SCX_ADDR])
	yy := y + int(mmap[SCY_ADDR])

	for i := 0; i < LCD_WIDTH; i++ {
		col := getTileColor(x%256, yy%256, mapAddr, tileAddr)
		setBufferPixel(x, y, col, pixels, pitch)
		x++
	}
}

// return the color of the pixel x, y for the sprite <index>
func getSpritePixel(x, y, pattern int) int {
	if x < 0 || y < 0 {
		return 0
	}
	return getPixel(uint16(0x8000+pattern*16), x, y)
}

// return the color of the pixel of the tile at address tileAddr and of
// coordinates x, y
func getPixel(tileAddr uint16, x, y int) int {
	addr := tileAddr + uint16(y*2)
	color := 2*GetBit(mmap[addr], uint8(7-x)) + GetBit(mmap[addr+1], uint8(7-x))
	return int(color)
}

// Return the color of the pixel at coordinates x, y
func getTileColor(x, y int, bgAddr, tileAddr uint16) int {
	// Get the tile corresponding to this coordinate
	tx := x / 8
	ty := y / 8

	// get the tile index in the tile map
	tIndexOffset := ty*32 + tx
	tIndex := mmap[bgAddr+uint16(tIndexOffset)]

	// Get the pixel x,y in the tile itself
	px := x % 8
	py := y % 8

	// get the tile address in the tile data table
	addr := tileAddr + uint16(tIndex*16)

	return getPixel(addr, px, py)
}

func clearScreen() {
	for i := 0; i < len(pixels); i++ {
		pixels[i] = 255
	}
}

func setLcdMode(mode int) {
	if mode == MODE_1_VBLANK {
		SetBit(mmap[STAT_ADDR], 1, 0)
		SetBit(mmap[STAT_ADDR], 0, 1)
	} else if mode == MODE_0_HBLANK {
		SetBit(mmap[STAT_ADDR], 1, 0)
		SetBit(mmap[STAT_ADDR], 0, 0)
	} else if mode == MODE_2_OAM_USED {
		SetBit(mmap[STAT_ADDR], 1, 1)
		SetBit(mmap[STAT_ADDR], 0, 0)
	} else if mode == MODE_3_OAM_VRAM_USED {
		SetBit(mmap[STAT_ADDR], 1, 1)
		SetBit(mmap[STAT_ADDR], 0, 1)
	}
}

func drawScanline(y int, lcdc uint8, pixels unsafe.Pointer, pitch int) {
	var tilecol, spritecol int

	bgAddr := getBackgroundTileMap()
	wdwAddr := getWindowTileMap()
	tileAddr := getTileDataTable()

	scx := int(mmap[SCX_ADDR])
	scy := int(mmap[SCY_ADDR])

	wx := int(mmap[WX_ADDR])
	wy := int(mmap[WY_ADDR])

	drawBgAndWindow := IsBitSet(lcdc, BG_WDW_ACTIVE)
	drawWindow := IsBitSet(lcdc, WDW_ACTIVE)
	drawSprites := IsBitSet(lcdc, SPRITE_ACTIVE)

	hblank := time.NewTicker(time.Microsecond * 90)

	for x := 0; x < LCD_WIDTH; x++ {
		if drawSprites {
			spritecol = getSpriteColor(x, y)
		}
		if drawBgAndWindow {
			tilecol = getTileColor((scx+x)%256, (scy+y)%256, bgAddr, tileAddr)

			if drawWindow && x >= wx && x <= wx+LCD_WIDTH && y >= wy && y <= wy+LCD_HEIGHT {
				tilecol = getTileColor(x, y, wdwAddr, tileAddr)
			}
		}
		if spritecol > tilecol {
			setBufferPixel(x, y, spritecol, pixels, pitch)
		} else {
			setBufferPixel(x, y, tilecol, pixels, pitch)
		}
	}
	setLcdMode(MODE_0_HBLANK)
	<-hblank.C
}

// Draw a single frame (144 lines + 10 "lines" of V-Blank (approx 1.1 ms))
func redraw() {
	lcdc := mem.GetLCDC()

	if !IsBitSet(lcdc, LCD_ACTIVE) {
		return
	}

	var pitch int
	var pixPtr unsafe.Pointer
	err := screenTex.Lock(nil, &pixPtr, &pitch)
	if err != nil {
		panic(err)
	}

	clearScreen()

	mem.SetLY(0x00)

	vblank := time.NewTicker(time.Millisecond * 16)
	scanline := time.NewTicker(time.Microsecond * 109)
	defer vblank.Stop()

	for y := 0; y < SCANLINES; y++ {
		if y < LCD_HEIGHT {
			drawScanline(y, lcdc, pixPtr, pitch)
		} else if y == LCD_HEIGHT {
			setLcdMode(MODE_1_VBLANK)
			cpu.RequestVBlankInterrupt()
		}
		mem.IncLY()
		<-scanline.C
	}

	screenTex.Update(nil, pixPtr, pitch)
	screenTex.Unlock()
	renderer.Clear()
	renderer.Copy(screenTex, nil, nil)
	renderer.Present()
	<-vblank.C
}
