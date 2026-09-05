package cmd

import (
	"fmt"
	"os"
	"strconv"

	"quickvm/internal/hyperv"
	"quickvm/internal/output"

	"github.com/spf13/cobra"
)

var (
	rdpCredentials string
	rdpCleanCreds  bool
)

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
for seamless login. Use --clean-creds to remove them when done.

Examples:
  quickvm rdp 1                               # RDP into VM 1
  quickvm rdp 2 -u admin                      # RDP with username
  quickvm rdp 1 -u "admin@password123"        # RDP with auto-login
  quickvm rdp 1 --clean-creds                 # Clean saved RDP credentials`,
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
			os.Exit(1)
		}

		// Get VM name for display
		vmName, err := manager.GetVMNameByIndex(cmd.Context(), index)
		if err != nil {
			output.PrintError("VM_GET_FAILED", "Failed to get VM", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to get VM: %v\n", err)
			}
			os.Exit(1)
		}

		// Get IP address first to show to user
		ip, err := manager.GetVMIPAddress(cmd.Context(), index)
		if err != nil {
			output.PrintError("IP_GET_FAILED", "Failed to get VM IP address", err.Error())
			if !output.IsJSON() {
				fmt.Printf("❌ Failed to get VM IP address: %v\n", err)
			}
			os.Exit(1)
		}

		// Handle credential cleanup request
		if rdpCleanCreds {
			if err := manager.DeleteRDPCredentials(cmd.Context(), ip); err != nil {
				output.PrintError("RDP_FAILED", "Failed to clean RDP credentials", err.Error())
				if !output.IsJSON() {
					fmt.Printf("❌ Failed to clean credentials: %v\n", err)
				}
				os.Exit(1)
			}
			if output.IsJSON() {
				output.PrintData(RDPResult{
					VMName:    vmName,
					VMIndex:   index,
					IPAddress: ip,
					Success:   true,
					Message:   "RDP credentials cleaned successfully",
				})
				return
			}
			fmt.Printf("✅ Saved RDP credentials for VM '%s' (%s) cleaned successfully!\n", vmName, ip)
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
			os.Exit(1)
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
			fmt.Printf("   - To clean up saved credentials, run: quickvm rdp %d --clean-creds\n", index)
		}
		fmt.Println("   - If connection fails, ensure Remote Desktop is enabled in the VM")
	},
}

func init() {
	rdpCmd.Flags().StringVarP(&rdpCredentials, "user", "u", "", "Credentials: \"username\" or \"username@password\"")
	rdpCmd.Flags().BoolVar(&rdpCleanCreds, "clean-creds", false, "Remove saved RDP credentials for this VM")
	rootCmd.AddCommand(rdpCmd)
}
