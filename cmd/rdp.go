package cmd

import (
	"fmt"
	"strconv"

	"quickvm/internal/hyperv"
	"quickvm/internal/output"

	"github.com/spf13/cobra"
)

var rdpCredentials string

var rdpCmd = &cobra.Command{
	Use:   "rdp <vm-index>",
	Short: "Open RDP connection to a VM",
	Long: `Open a Remote Desktop connection to a Hyper-V virtual machine.

This command gets the VM's IP address and opens the Windows Remote Desktop client (mstsc.exe).

Requirements:
  - VM must be running
  - VM must have integration services installed
  - VM must have an IPv4 address assigned
  - Remote Desktop must be enabled in the VM

Credentials format:
  - Username only: -u "username"
  - Username with password: -u "username@password"
  - Domain user: -u "domain\username@password"

When password is provided, credentials are saved to Windows Credential Manager
for seamless login.

Examples:
  quickvm rdp 1                               # RDP into VM 1
  quickvm rdp 2 -u admin                      # RDP with username
  quickvm rdp 1 -u "admin@password123"        # RDP with auto-login`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := hyperv.NewManager()

		// Parse VM index
		index, err := strconv.Atoi(args[0])
		if err != nil {
			output.PrintError("INVALID_INDEX", "Invalid VM index", args[0])
			if !output.IsJSON() {
				fmt.Printf("❌ Invalid VM index: %s\n", args[0])
			}
			return
		}

		// Get VM name for display
		vmName, err := manager.GetVMNameByIndex(cmd.Context(), index)
		if err != nil {
			output.PrintError("VM_GET_FAILED", "Failed to get VM", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to get VM: %v\n", err)
			}
			return
		}

		// Get IP address first to show to user
		ip, err := manager.GetVMIPAddress(cmd.Context(), index)
		if err != nil {
			output.PrintError("IP_GET_FAILED", "Failed to get VM IP address", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to get VM IP address: %v\n", err)
			}
			return
		}

		if !output.IsJSON() {
			fmt.Printf("🔗 Connecting to VM '%s' at %s...\n", vmName, ip)
		}

		// Parse credentials for display
		creds := hyperv.ParseCredentials(rdpCredentials)
		if creds.Password != "" && !output.IsJSON() {
			fmt.Println("🔐 Saving credentials to Windows Credential Manager...")
		}

		if err := manager.ConnectRDPByIP(cmd.Context(), ip, rdpCredentials); err != nil {
			output.PrintError("RDP_FAILED", "Failed to open RDP", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to open RDP: %v\n", err)
			}
			return
		}

		// JSON output for AI agents
		if output.IsJSON() {
			output.PrintData(RDPResult{
				VMName:    vmName,
				VMIndex:   index,
				IPAddress: ip,
				Success:   true,
				Message:   "RDP client opened successfully",
			})
			return
		}

		fmt.Println("✅ RDP client opened successfully!")
		fmt.Println()
		fmt.Println("💡 Tips:")
		fmt.Printf("   - IP address: %s\n", ip)
		if creds.Username != "" {
			fmt.Printf("   - Username: %s\n", creds.Username)
		}
		if creds.Password != "" {
			fmt.Println("   - Credentials saved for auto-login")
		}
		fmt.Println("   - If connection fails, ensure Remote Desktop is enabled in the VM")
	},
}

func init() {
	rdpCmd.Flags().StringVarP(&rdpCredentials, "user", "u", "", "Credentials: \"username\" or \"username@password\"")
	rootCmd.AddCommand(rdpCmd)
}
