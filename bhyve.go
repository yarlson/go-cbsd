package cbsd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	ActionStart  = "bstart"
	ActionStop   = "bstop"
	ActionRemove = "bremove"
	ActionCreate = "bcreate"

	NoSuchDomainError = "No such domain"
)

type BHyve struct {
	JName    string
	JID      int
	VmRam    int
	VmCPUs   int
	VmOSType string
	Status   string
	VNC      string
}

// BHyveCreate is a struct that contains all fields to run bcreate command
// every none empty/none nil field will be converted to name=value pair, the name will be taken form a json tag
type BHyveCreate struct {
	Name               string `json:"jname,omitempty"`         // A short jail name
	XHCI               string `json:"xhci,omitempty"`          // eXtensible Host Controller Interface (xHCI) USB controller
	AutoStart          *bool  `json:"astart,omitempty"`        // Automatically start Jail when system boot
	RelativePath       *bool  `json:"relative_path,omitempty"` // Use relative path
	Path               string `json:"path,omitempty"`          // Jail path
	Data               string `json:"data,omitempty"`          // Data path
	RCConf             string `json:"rcconf,omitempty"`        // rc.conf path
	Hostname           string `json:"host_hostname,omitempty"` // Full (FQDN) jail hostname
	IP4Addr            string `json:"ip4_addr,omitempty"`      // IPv4 address
	NicHWAddr          string `json:"nic_hwaddr,omitempty"`    // MAC address
	ZfsSnapSrc         string `json:"zfs_snapsrc,omitempty"`   // Use this ZFS snapshot as source for jail data
	RunASAP            *bool  `json:"runasap,omitempty"`       // Start jail ASAP upon creation
	Interface          string `json:"interface,omitempty"`
	RCtlNice           string `json:"rctl_nice,omitempty"`
	Emulator           string `json:"emulator,omitempty"`
	ImgSize            string `json:"imgsize,omitempty"`
	ImgType            string `json:"imgtype,omitempty"`
	VmCPUs             string `json:"vm_cpus,omitempty"`
	VmRAM              string `json:"vm_ram,omitempty"`
	VmOSType           string `json:"vm_os_type,omitempty"`
	VmEFI              string `json:"vm_efi,omitempty"`
	IsoSite            string `json:"iso_site,omitempty"`
	IsoImg             string `json:"iso_img,omitempty"`
	RegisterIsoName    string `json:"register_iso_name,omitempty"`
	RegisterIsoAs      string `json:"register_iso_as,omitempty"`
	VmHostBridge       string `json:"vm_hostbridge,omitempty"`
	BhyveFlags         string `json:"bhyve_flags,omitempty"`
	VirtioType         string `json:"virtio_type,omitempty"`
	VmOSProfile        string `json:"vm_os_profile,omitempty"`
	SwapSize           string `json:"swapsize,omitempty"`
	VmIsoPath          string `json:"vm_iso_path,omitempty"`
	VmGuestFS          string `json:"vm_guestfs,omitempty"`
	VmVNCPort          string `json:"vm_vnc_port,omitempty"`
	BhyveGenerateAcpi  string `json:"bhyve_generate_acpi,omitempty"`
	BhyveWireMemory    string `json:"bhyve_wire_memory,omitempty"`
	BhyveRtsKeepsUtc   string `json:"bhyve_rts_keeps_utc,omitempty"`
	BhyveForceMsiIrq   string `json:"bhyve_force_msi_irq,omitempty"`
	BhyveX2ApicMode    string `json:"bhyve_x2apic_mode,omitempty"`
	BhyveMpTableGen    string `json:"bhyve_mptable_gen,omitempty"`
	BhyveIgnoreMsrAcc  string `json:"bhyve_ignore_msr_acc,omitempty"`
	CdVncWait          string `json:"cd_vnc_wait,omitempty"`
	BhyveVNCResolution string `json:"bhyve_vnc_resolution,omitempty"`
	BhyveVNCTcpBind    string `json:"bhyve_vnc_tcp_bind,omitempty"`
	BhyveVNCVgaConf    string `json:"bhyve_vnc_vgaconf,omitempty"`
	NicDriver          string `json:"nic_driver,omitempty"`
	VNCPassword        string `json:"vnc_password,omitempty"`
	MediaAutoEject     string `json:"media_auto_eject,omitempty"`
	VmCPUTopology      string `json:"vm_cpu_topology,omitempty"`
	DebugEngine        string `json:"debug_engine,omitempty"`
	CdBootFirmware     string `json:"cd_boot_firmware,omitempty"`
	Jailed             string `json:"jailed,omitempty"`
	OnPowerOff         string `json:"on_poweroff,omitempty"`
	OnReboot           string `json:"on_reboot,omitempty"`
	OnCrash            string `json:"on_crash,omitempty"`
}

type BHyveService struct {
	exec Exec
}

func (b *BHyveService) do(ctx context.Context, action, instanceId string) ([]byte, error) {
	b.exec.SetEnv("NOCOLOR", "1")
	return b.exec.Command(ctx, "cbsd", action, "inter=0", fmt.Sprintf("jname=%s", instanceId))
}

func (b *BHyveService) Start(ctx context.Context, instanceId string) error {
	bt, err := b.do(ctx, ActionStart, instanceId)
	if err != nil {
		return err
	}

	output := string(bt)
	if strings.Contains(output, NoSuchDomainError) {
		return errors.New(output)
	}

	return nil
}

func (b *BHyveService) Stop(ctx context.Context, instanceId string) error {
	bt, err := b.do(ctx, ActionStop, instanceId)
	if err != nil {
		return err
	}

	output := string(bt)
	if strings.Contains(output, NoSuchDomainError) {
		return errors.New(output)
	}

	return nil
}

func (b *BHyveService) Remove(ctx context.Context, instanceId string) error {
	bt, err := b.do(ctx, ActionRemove, instanceId)
	if err != nil {
		return err
	}

	output := string(bt)
	if !strings.Contains(output, NoSuchDomainError) {
		return nil
	}

	lines := strings.Split(string(output), "\n")
	var line string
	for _, line = range lines {
		if strings.Contains(line, NoSuchDomainError) {
			break
		}
	}

	return errors.New(line)
}

func (b *BHyveService) List(ctx context.Context) ([]BHyve, error) {
	b.exec.SetEnv("NOCOLOR", "1")
	output, err := b.exec.Command(ctx, "cbsd", "bls", "header=0", "display=jname,jid,vm_ram,vm_cpus,vm_os_type,status,vnc_port")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")

	var bHyves []BHyve
	for _, line := range lines {
		if len(line) <= 2 {
			continue
		}
		fields := strings.Fields(line)
		ima := BHyve{}
		ima.JName = fields[0]
		ima.JID, _ = strconv.Atoi(fields[1])
		ima.VmRam, _ = strconv.Atoi(fields[2])
		ima.VmCPUs, _ = strconv.Atoi(fields[3])
		ima.VmOSType = fields[4]
		ima.Status = fields[5]
		ima.VNC = fields[6]
		bHyves = append(bHyves, ima)
	}

	return bHyves, nil
}

func (b *BHyveService) Create(ctx context.Context, createData *BHyveCreate) ([]byte, error) {
	b.exec.SetEnv("NOCOLOR", "1")
	return b.exec.CommandWithInterface(ctx, "cbsd", createData, ActionCreate, "inter=0")
}
