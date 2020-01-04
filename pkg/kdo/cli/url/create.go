package url

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/redhat-developer/odo-fork/pkg/config"
	"github.com/redhat-developer/odo-fork/pkg/kdo/genericclioptions"
	"github.com/redhat-developer/odo-fork/pkg/log"
	"github.com/redhat-developer/odo-fork/pkg/url"

	"github.com/redhat-developer/odo-fork/pkg/util"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const createRecommendedCommandName = "create"

var (
	urlCreateShortDesc = `Create a URL for a component`
	urlCreateLongDesc  = ktemplates.LongDesc(`Create a URL for a component.

	The created URL can be used to access the specified component from outside the cluster.
	`)
	urlCreateExample = ktemplates.Examples(`

  	# Create a URL with a specific port
	%[1]s example.<ingress domain> --port 8080
  
	# Create a URL by automatic detection of port (only for components which expose only one service port) 
	%[1]s example.<ingress domain>
  
	  `)
)

// URLCreateOptions encapsulates the options for the odo url create command
type URLCreateOptions struct {
	localConfigInfo  *config.LocalConfigInfo
	componentContext string
	urlName          string
	urlPort          int
	ingressDomain    string
	componentPort    int
	*genericclioptions.Context
}

// NewURLCreateOptions creates a new UrlCreateOptions instance
func NewURLCreateOptions() *URLCreateOptions {
	return &URLCreateOptions{}
}

// Complete completes UrlCreateOptions after they've been Created
func (o *URLCreateOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	o.Context = genericclioptions.NewContext(cmd)
	o.componentPort, err = url.GetValidPortNumber(o.Client, o.urlPort, o.Component(), o.Application)
	o.urlPort = o.componentPort
	if err != nil {
		return err
	}
	// o.urlName = url.GetURLName(o.Component(), o.componentPort)
	o.urlName = o.Component()
	if len(args) == 0 {
		return fmt.Errorf("An ingress domain must be provided")
	} else {
		o.ingressDomain = args[0]
	}
	o.localConfigInfo, err = config.NewLocalConfigInfo(o.componentContext)

	return
}

// Validate validates the UrlCreateOptions based on completed values
func (o *URLCreateOptions) Validate() (err error) {

	// Check if exist
	for _, localUrl := range o.localConfigInfo.GetUrl() {
		if o.urlName == localUrl.Name {
			return fmt.Errorf("the url %s already exists in the application: %s", o.urlName, o.Application)
		}
	}

	// Check if url name is more than 63 characters long
	if len(o.urlName) > 63 {
		return fmt.Errorf("url name must be shorter than 63 characters")
	}

	if !util.CheckOutputFlag(o.OutputFlag) {
		return fmt.Errorf("given output format %s is not supported", o.OutputFlag)
	}

	return
}

// Run contains the logic for the odo url create command
func (o *URLCreateOptions) Run() (err error) {
	err = o.localConfigInfo.SetConfiguration("url", config.ConfigUrl{Name: o.urlName, Port: o.urlPort, Host: o.ingressDomain})
	if err != nil {
		return errors.Wrapf(err, "failed to persist the component settings to config file")
	}
	log.Successf("URL %s created", o.urlName)
	log.Infof("\nRun `udo push` to apply URL: %s", o.urlName)
	return
}

// NewCmdURLCreate implements the odo url create command.
func NewCmdURLCreate(name, fullName string) *cobra.Command {
	o := NewURLCreateOptions()
	urlCreateCmd := &cobra.Command{
		Use:     name + " [ingress domain]",
		Short:   urlCreateShortDesc,
		Long:    urlCreateLongDesc,
		Example: fmt.Sprintf(urlCreateExample, fullName),
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}
	urlCreateCmd.Flags().IntVarP(&o.urlPort, "port", "", -1, "port number for the url of the component, required in case of components which expose more than one service port")
	// _ = urlCreateCmd.MarkFlagRequired("port")
	genericclioptions.AddOutputFlag(urlCreateCmd)
	genericclioptions.AddContextFlag(urlCreateCmd, &o.componentContext)
	// completion.RegisterCommandFlagHandler(urlCreateCmd, "context", completion.FileCompletionHandler)
	return urlCreateCmd
}
