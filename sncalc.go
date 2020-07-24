/*
Program Name    : sncalc.go
Binary Name		: sncalc
Build Command	: go build -o sncalc

Program Purpose : Subnet Calculator 
                  Based on subnetting chapters in below books:
                    "CCNA Routing and Switching Study Guide" by Todd Lammle - 2013 edition
                  Similar results obtained from:
                    calculator.net/ip-subnet-calculator.html
                    tunnelsup.com/subnet-calculator
					
Description     : This tool takes an network IP address and CIDR value as input and provides useful 
                  information like the network address, broadcast address, host range in the network, 
                  and wildcard mask, among others. This website also provides a list of subnets possible 
                  with the IP & CIDR provided so that a large network can be subdivided into smaller 
                  manageable subnets.
				  
Date            : 14-MAY-2020
Author          : Sam

*/


package main

import (
	"fmt"
	"strings"
	"strconv"
	"math"
	"errors"
	"os"
)


const (
	minfo_version string = "sncalc 1.0"
	copyright string = "Copyright (C) 2020 Sanjeev Medhi."
)

const (
	ipTotalBitCount int = 32
	maxOctetDecimal int = 255
	maxNetworkBitsForUsefulHosts int = 30
)


var (
	ipv4 string = "192.168.1.0"
	cidr int = 26
	
	//inputSubnetIp string = "172.16.0.0"
	inputSubnetCidr int = 24
	//requiredHostAddressesPerSubnet int = 30
	
	metricMap = make(map[string]string, 0)
    //metricMap = map[string]string{}
    
    subnetListSlice = make([]string, 0)
)


func main() {

	ipValidationMap, err := ipValidation(ipv4)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		for k, v := range ipValidationMap {
			metricMap[k] = v
		}
	}
	
	subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", "Network Address", "Usable Host Range", "Broadcast Address"))
	subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", "---------------", "-----------------", "-----------------"))
	
	cidrToSubnetMaskMap, err := cidrToSubnetMask(cidr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		subnetMask := cidrToSubnetMaskMap["Subnet Mask"]
		subnetCalcMap, err := subnetCalc(ipv4, cidr, subnetMask)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			for k, v := range subnetCalcMap {
				metricMap[k] = v
			}
		}
		
		for k, v := range cidrToSubnetMaskMap {
			metricMap[k] = v
		}
	}
	
	hostsPerSubnetCalcMap, err := hostsPerSubnetCalc()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		for k, v := range hostsPerSubnetCalcMap {
			metricMap[k] = v
		}
	}
	
	fmt.Printf("\n")
	metricMapDisplay()
	fmt.Printf("\n")
	
	fmt.Printf("\n")
	fmt.Printf("Number of Subnets: %s\n", metricMap["Number of Subnets"])
	fmt.Printf("%s", metricMap["Subnet List"])
	for _, v := range subnetListSlice {
		fmt.Printf("%v\n", v)
	}
	fmt.Printf("\n")
	
	
	
	// SEGMENT NOT READY
	/* 
	fmt.Printf("Provided subnet CIDR: %v \n", inputSubnetCidr)
	fmt.Printf("Required number of hosts per subnet: %v \n", requiredHostAddressesPerSubnet)
	var s []int
	var availableHostBits int
	var minNeededHostsPerSubnet int
	hostBits := ipTotalBitCount - inputSubnetCidr
	fmt.Printf("Host bits available for subnetting: %v \n", hostBits)
	for i := 0; i <= hostBits; i++ {
		hostCnt := math.Pow(2, float64(i)) - float64(2)
		s = append(s, int(hostCnt))
	}
	
	for i := 1; i <= hostBits; i++ {
		if requiredHostAddressesPerSubnet > s[i-1] && requiredHostAddressesPerSubnet <= s[i] {
			availableHostBits = i
			minNeededHostsPerSubnet = s[i]
		}
	}
	fmt.Printf("Host bits available after subnetting: %v \n", availableHostBits)
	fmt.Printf("Host count that will fullfil the requirement of %v hosts per subnet: %v \n", requiredHostAddressesPerSubnet, minNeededHostsPerSubnet)
	hostBitsForSubnetting := hostBits - availableHostBits
	subnetCnt := int(math.Pow(2, float64(hostBitsForSubnetting)))
	fmt.Printf("Number of possible subnets: %v \n", subnetCnt)
	newSubnetCidr := inputSubnetCidr + (hostBits - availableHostBits)
	fmt.Printf("New Subnet Mask: %v \n", newSubnetCidr)
	
	map1, err := cidrToSubnetMask(newSubnetCidr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		newSubnetMask := map1["Subnet Mask"]
		s := strings.Split(newSubnetMask, ".")
	    s1, _ := strconv.Atoi(s[0])
		s2, _ := strconv.Atoi(s[1])
		s3, _ := strconv.Atoi(s[2])
		s4, _ := strconv.Atoi(s[3])
		
		var subnetNbr int
		if newSubnetCidr >= 24 && newSubnetCidr <= maxNetworkBitsForUsefulHosts {
			subnetNbr = s4
		} else if newSubnetCidr >= 16 && newSubnetCidr < 24 {
			subnetNbr = s3
		} else if newSubnetCidr >= 8 && newSubnetCidr < 16 {
			subnetNbr = s2
		} else if newSubnetCidr >= 1 && newSubnetCidr < 8 {
			subnetNbr = s1	
		}
		fmt.Printf("Subnet Multiplier: %v \n", 256 - subnetNbr)
	}
	*/
	
}


func metricMapDisplay() {
	fmt.Printf("%-40s: %s\n", "IP Address", metricMap["IP Address"])
	fmt.Printf("%-40s: %s\n", "Network Address", metricMap["Network Address"])
	fmt.Printf("%-40s: %s\n", "Usable Host IP Range", metricMap["Usable Host IP Range"])
	fmt.Printf("%-40s: %s\n", "Broadcast Address", metricMap["Broadcast Address"])
	fmt.Printf("%-40s: %s\n", "Total Hosts per Subnet", metricMap["Total Hosts per Subnet"])
	fmt.Printf("%-40s: %s\n", "Usable Hosts per Subnet", metricMap["Usable Hosts per Subnet"])
	fmt.Printf("%-40s: %s\n", "Subnet Mask", metricMap["Subnet Mask"])
	fmt.Printf("%-40s: %s\n", "Wildcard Mask", metricMap["Wildcard Mask"])
	fmt.Printf("%-40s: %s\n", "Binary Subnet Mask", metricMap["Binary Subnet Mask"])
	fmt.Printf("%-40s: %s\n", "CIDR Notation", metricMap["CIDR Notation"])
	fmt.Printf("%-40s: %s\n", "Binary Octets", metricMap["Binary Octets"])
	fmt.Printf("%-40s: %s\n", "Network Bits (total masked bits)", metricMap["Network Bits (total masked bits)"])
	fmt.Printf("%-40s: %s\n", "Hosts Bits (unmasked bits)", metricMap["Hosts Bits (unmasked bits)"])
	
}


func ipValidation(ipv4 string) (map[string]string, error) {
	var ipValidationMap = make(map[string]string, 0)

	ipValidationMap["IP Address"] = ipv4

	ipSlice := strings.Split(ipv4, ".")

	octet1Decimal, _ := strconv.Atoi(ipSlice[0])
	octet2Decimal, _ := strconv.Atoi(ipSlice[1])
	octet3Decimal, _ := strconv.Atoi(ipSlice[2])
	octet4Decimal, _ := strconv.Atoi(ipSlice[3])
	
	if (octet1Decimal > maxOctetDecimal) || (octet2Decimal > maxOctetDecimal) || (octet3Decimal > maxOctetDecimal) || (octet4Decimal > maxOctetDecimal) {
		myErr := errors.New("ERROR: Max octet decimal value can be 255.")
		return nil, myErr
	}
	
	octet1 := fmt.Sprintf("%08b", octet1Decimal)
	octet2 := fmt.Sprintf("%08b", octet2Decimal)
	octet3 := fmt.Sprintf("%08b", octet3Decimal)
	octet4 := fmt.Sprintf("%08b", octet4Decimal)			
	
	binaryOctets := fmt.Sprintf("%v.%v.%v.%v", octet1, octet2, octet3, octet4)
	
	ipValidationMap["Binary Octets"] = binaryOctets
	
	return ipValidationMap, nil
}


func hostsPerSubnetCalc() (map[string]string, error) {
	var hostsPerSubnetCalcMap = make(map[string]string, 0)
	if cidr <= maxNetworkBitsForUsefulHosts {
		unmaskedBits := ipTotalBitCount - cidr
		hostsPerSubnet := math.Pow(2, float64(unmaskedBits))
		hostsPerSubnetCalcMap["Total Hosts per Subnet"] = fmt.Sprintf("%v   (2^unmasked bits) => (2^%v)", int(hostsPerSubnet), unmaskedBits)
		usableHostsPerSubnet := math.Pow(2, float64(unmaskedBits)) - 2
		hostsPerSubnetCalcMap["Usable Hosts per Subnet"] = fmt.Sprintf("%v   (2^unmasked bits - 2) => (2^%v - 2)", int(usableHostsPerSubnet), unmaskedBits)
	} else if cidr > maxNetworkBitsForUsefulHosts {
		unmaskedBits := ipTotalBitCount - cidr
		hostsPerSubnet := math.Pow(2, float64(unmaskedBits))
		hostsPerSubnetCalcMap["Total Hosts per Subnet"] = fmt.Sprintf("%v   (2^unmasked bits) => (2^%v)", int(hostsPerSubnet), unmaskedBits)
		hostsPerSubnetCalcMap["Usable Hosts per Subnet"] = fmt.Sprintf("%v", 0)
	} else {
		hostsPerSubnetCalcMap["Total Hosts per Subnet"] = fmt.Sprintf("%v", 0)
		hostsPerSubnetCalcMap["Usable Hosts per Subnet"] = fmt.Sprintf("%v", 0)
	}
	
	return hostsPerSubnetCalcMap, nil
}


func subnetCalc(ipv4 string, cidr int, subnetMask string) (map[string]string, error) {
	var subnetCalcMap = make(map[string]string, 0)
	ip1 := strings.Split(ipv4, ".")
	
	s := strings.Split(subnetMask, ".")
	s1, _ := strconv.Atoi(s[0])
	s2, _ := strconv.Atoi(s[1])
	s3, _ := strconv.Atoi(s[2])
	s4, _ := strconv.Atoi(s[3])
	
	if cidr >= 24 && cidr <= maxNetworkBitsForUsefulHosts {
		maskedBits := cidr - 24
		numberOfSubnets := math.Pow(2, float64(maskedBits))
		subnetCalcMap["Number of Subnets"] = fmt.Sprintf("%v   (2^masked bits on 4th octet) => (2^%v)", numberOfSubnets, maskedBits)
		networkPortion := fmt.Sprintf("%v.%v.%v", ip1[0], ip1[1], ip1[2])
		subnetNbr := s4
		octetPosition := 4
		subnetCalcMap["Subnet List"] = fmt.Sprintf("All %v of the Possible /%v Networks for %v.%v.%v.* (valid subnets at 4th octet):\n", numberOfSubnets, cidr, ip1[0], ip1[1], ip1[2])
		subnetListMap, err := subnetList(networkPortion, subnetNbr, octetPosition, ip1[3])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			for k, v := range subnetListMap {
				subnetCalcMap[k] = v
			}
		}
		
	} else if cidr >= 16 && cidr < 24 {
		maskedBits := cidr - 16
		numberOfSubnets := math.Pow(2, float64(maskedBits))
		subnetCalcMap["Number of Subnets"] = fmt.Sprintf("%v   (2^masked bits on 3rd octet) => (2^%v)", numberOfSubnets, maskedBits)
		networkPortion := fmt.Sprintf("%v.%v", ip1[0], ip1[1])
		subnetNbr := s3
		octetPosition := 3
		subnetCalcMap["Subnet List"] = fmt.Sprintf("All %v of the Possible /%v Networks for %v.%v.*.* (valid subnets at 3rd octet):\n", numberOfSubnets, cidr, ip1[0], ip1[1])
		subnetListMap, err := subnetList(networkPortion, subnetNbr, octetPosition, ip1[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			for k, v := range subnetListMap {
				subnetCalcMap[k] = v
			}
		}
		
	} else if cidr >= 8 && cidr < 16 {
		maskedBits := cidr - 8
		numberOfSubnets := math.Pow(2, float64(maskedBits))
		subnetCalcMap["Number of Subnets"] = fmt.Sprintf("%v   (2^masked bits on 2nd octet) => (2^%v)", numberOfSubnets, maskedBits)
		networkPortion := fmt.Sprintf("%v", ip1[0])
		subnetNbr := s2
		octetPosition := 2
		subnetCalcMap["Subnet List"] = fmt.Sprintf("All %v of the Possible /%v Networks for %v.*.*.* (valid subnets at 2nd octet):\n", numberOfSubnets, cidr, ip1[0])
		subnetListMap, err := subnetList(networkPortion, subnetNbr, octetPosition, ip1[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			for k, v := range subnetListMap {
				subnetCalcMap[k] = v
			}
		}
		
	} else if cidr >= 1 && cidr < 8 {
		maskedBits := cidr - 0
		numberOfSubnets := math.Pow(2, float64(maskedBits))
		subnetCalcMap["Number of Subnets"] = fmt.Sprintf("%v   (2^masked bits on 1st octet) => (2^%v)", numberOfSubnets, maskedBits)
		networkPortion := fmt.Sprintf("")
		subnetNbr := s1
		octetPosition := 1
		subnetCalcMap["Subnet List"] = fmt.Sprintf("All %v of the Possible /%v Networks (valid subnets at 1st octet):\n", numberOfSubnets, cidr)
		subnetListMap, err := subnetList(networkPortion, subnetNbr, octetPosition, ip1[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			for k, v := range subnetListMap {
				subnetCalcMap[k] = v
			}
		}
		
	} else {
		subnetCalcMap["Number of Subnets"] = fmt.Sprintf("%v", 0)
	}
	
	return subnetCalcMap, nil
}


func subnetList(networkPortion string, subnetNbr int, octetPosition int, octetValue string) (map[string]string, error) {
	var subnetListMap = make(map[string]string, 0)
	
	blockSize := 256 - subnetNbr
	i := blockSize
	octetValueInt , _ := strconv.Atoi(octetValue)
	
	if octetPosition == 4 {
		nAddr := fmt.Sprintf("%v.%v", networkPortion, 0)
		startIP := fmt.Sprintf("%v.%v", networkPortion, 1)
		endIP := fmt.Sprintf("%v.%v", networkPortion, blockSize - 2)
		usableHostIPRange := fmt.Sprintf("%v - %v", startIP, endIP)
		broadcastAddress := fmt.Sprintf("%v.%v", networkPortion, blockSize - 1)
		
		if octetValueInt >= 0 && octetValueInt < i {
			subnetListMap["Network Address"] = nAddr
			subnetListMap["Usable Host IP Range"] = usableHostIPRange
			subnetListMap["Broadcast Address"] = broadcastAddress
			subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v [current]", nAddr, usableHostIPRange, broadcastAddress))
			
		} else {
			subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", nAddr, usableHostIPRange, broadcastAddress))
		}
		
		for ; i <= subnetNbr; i = i + blockSize {
			nAddr := fmt.Sprintf("%v.%v", networkPortion, i)
			startIP := fmt.Sprintf("%v.%v", networkPortion, i + 1)
			endIP := fmt.Sprintf("%v.%v", networkPortion, i + blockSize - 2)
			usableHostIPRange := fmt.Sprintf("%v - %v", startIP, endIP)
			broadcastAddress := fmt.Sprintf("%v.%v", networkPortion, i + blockSize - 1)
			
			if octetValueInt >= i && octetValueInt < i + blockSize {				
				subnetListMap["Network Address"] = nAddr
				subnetListMap["Usable Host IP Range"] = usableHostIPRange
				subnetListMap["Broadcast Address"] = broadcastAddress
				subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v [current]", nAddr, usableHostIPRange, broadcastAddress))
				
			} else {
				subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", nAddr, usableHostIPRange, broadcastAddress))
			}
		}
	}
	
	if octetPosition == 3 {
		nAddr := fmt.Sprintf("%v.%v.%v", networkPortion, 0, 0)
		startIP := fmt.Sprintf("%v.%v.%v", networkPortion, 0, 1)
		endIP := fmt.Sprintf("%v.%v.%v", networkPortion, i - 1, 254)
		usableHostIPRange := fmt.Sprintf("%v - %v", startIP, endIP)
		broadcastAddress := fmt.Sprintf("%v.%v.%v", networkPortion, i - 1, maxOctetDecimal)
		
		if octetValueInt >= 0 && octetValueInt < i {
			subnetListMap["Network Address"] = nAddr
			subnetListMap["Usable Host IP Range"] = usableHostIPRange
			subnetListMap["Broadcast Address"] = broadcastAddress
			subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v [current]", nAddr, usableHostIPRange, broadcastAddress))
		} else {
			subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", nAddr, usableHostIPRange, broadcastAddress))
		}
		
		for ; i <= subnetNbr; i = i + blockSize {
			nAddr := fmt.Sprintf("%v.%v.%v", networkPortion, i, 0)
			startIP := fmt.Sprintf("%v.%v.%v", networkPortion, i, 1)
			endIP := fmt.Sprintf("%v.%v.%v", networkPortion, i + blockSize - 1, 254)
			usableHostIPRange := fmt.Sprintf("%v - %v", startIP, endIP)
			broadcastAddress := fmt.Sprintf("%v.%v.%v", networkPortion, i + blockSize - 1, maxOctetDecimal)
			
			if octetValueInt >= i && octetValueInt < i + blockSize {
				subnetListMap["Network Address"] = nAddr
				subnetListMap["Usable Host IP Range"] = usableHostIPRange
				subnetListMap["Broadcast Address"] = broadcastAddress
				subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v [current]", nAddr, usableHostIPRange, broadcastAddress))
			} else {
				subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", nAddr, usableHostIPRange, broadcastAddress))
			}
		}
	}
	
	if octetPosition == 2 {
		nAddr := fmt.Sprintf("%v.%v.%v.%v", networkPortion, 0, 0, 0)
		startIP := fmt.Sprintf("%v.%v.%v.%v", networkPortion, 0, 0, 1)
		endIP := fmt.Sprintf("%v.%v.%v.%v", networkPortion, i - 1, maxOctetDecimal, 254)
		usableHostIPRange := fmt.Sprintf("%v - %v", startIP, endIP)
		broadcastAddress := fmt.Sprintf("%v.%v.%v.%v", networkPortion, i - 1, maxOctetDecimal, maxOctetDecimal)
		
		if octetValueInt >= 0 && octetValueInt < i {
			subnetListMap["Network Address"] = nAddr
			subnetListMap["Usable Host IP Range"] = usableHostIPRange
			subnetListMap["Broadcast Address"] = broadcastAddress
			subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v [current]", nAddr, usableHostIPRange, broadcastAddress))
		} else {
			subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", nAddr, usableHostIPRange, broadcastAddress))
		}
		
		for ; i <= subnetNbr; i = i + blockSize {
			nAddr := fmt.Sprintf("%v.%v.%v.%v", networkPortion, i, 0, 0)
			startIP := fmt.Sprintf("%v.%v.%v.%v", networkPortion, i, 0, 1)
			endIP := fmt.Sprintf("%v.%v.%v.%v", networkPortion, i + blockSize - 1, maxOctetDecimal, 254)
			usableHostIPRange := fmt.Sprintf("%v - %v", startIP, endIP)
			broadcastAddress := fmt.Sprintf("%v.%v.%v.%v", networkPortion, i + blockSize - 1, maxOctetDecimal, maxOctetDecimal)
			
			if octetValueInt >= i && octetValueInt < i + blockSize {
				subnetListMap["Network Address"] = nAddr
				subnetListMap["Usable Host IP Range"] = usableHostIPRange
				subnetListMap["Broadcast Address"] = broadcastAddress
				subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v [current]", nAddr, usableHostIPRange, broadcastAddress))
			} else {
				subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", nAddr, usableHostIPRange, broadcastAddress))
			}
		}
	}
	
	if octetPosition == 1 {
		nAddr := fmt.Sprintf("%v.%v.%v.%v", 0, 0, 0, 0)
		startIP := fmt.Sprintf("%v.%v.%v.%v", 0, 0, 0, 1)
		endIP := fmt.Sprintf("%v.%v.%v.%v", blockSize - 1, maxOctetDecimal, maxOctetDecimal, 254)
		usableHostIPRange := fmt.Sprintf("%v - %v", startIP, endIP)
		broadcastAddress := fmt.Sprintf("%v.%v.%v.%v", blockSize - 1, maxOctetDecimal, maxOctetDecimal, maxOctetDecimal)
		
		if octetValueInt >= 0 && octetValueInt < i {
			subnetListMap["Network Address"] = nAddr
			subnetListMap["Usable Host IP Range"] = usableHostIPRange
			subnetListMap["Broadcast Address"] = broadcastAddress
			subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v [current]", nAddr, usableHostIPRange, broadcastAddress))
		} else {
			subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", nAddr, usableHostIPRange, broadcastAddress))
		}
		
		for ; i <= subnetNbr; i = i + blockSize {
			nAddr := fmt.Sprintf("%v.%v.%v.%v", i, 0, 0, 0)
			startIP := fmt.Sprintf("%v.%v.%v.%v", i, 0, 0, 1)
			endIP := fmt.Sprintf("%v.%v.%v.%v", i + blockSize - 1, maxOctetDecimal, maxOctetDecimal, 254)
			usableHostIPRange := fmt.Sprintf("%v - %v", startIP, endIP)
			broadcastAddress := fmt.Sprintf("%v.%v.%v.%v", i + blockSize - 1, maxOctetDecimal, maxOctetDecimal, maxOctetDecimal)
			
			if octetValueInt >= i && octetValueInt < i + blockSize {
				subnetListMap["Network Address"] = nAddr
				subnetListMap["Usable Host IP Range"] = usableHostIPRange
				subnetListMap["Broadcast Address"] = broadcastAddress
				subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v [current]", nAddr, usableHostIPRange, broadcastAddress))
			} else {
				subnetListSlice = append(subnetListSlice, fmt.Sprintf("  %-20v %-40v %v", nAddr, usableHostIPRange, broadcastAddress))
			}
		}
	}
	
	return subnetListMap, nil
}



func cidrToSubnetMask(cidr int) (map[string]string, error) {
	var cidrToSubnetMaskMap = make(map[string]string, 0)

	cidrToSubnetMaskMap["CIDR Notation"] = fmt.Sprintf("/%v", cidr)
	
	v_cidr := cidr
	
	if v_cidr > ipTotalBitCount {
		myErr := errors.New("ERROR: Max network mask (bits) can be 32")
		return nil, myErr
	}
	
	str1 := ""
	cnt1 := 0
	cnt2 := 1
	cnt3 := 1
	
	for i := 1; i <= ipTotalBitCount; i++ {
	    cnt1++
		if cnt1 == 8 {
		  cnt2++
		  if cnt2 == 5 {
			if v_cidr > 0 {
				str1 = fmt.Sprintf(str1+"1")
			} else {
				str1 = fmt.Sprintf(str1+"0")
			}
			break
		  } else {
			if v_cidr > 0 {
				str1 = fmt.Sprintf(str1+"1.")
			} else {
				str1 = fmt.Sprintf(str1+"0.")
			}
			cnt3++
		  }
		  cnt1 = 0
		} else {
		  if v_cidr > 0 {
		  	str1 = fmt.Sprintf(str1+"1")
		  } else {
			str1 = fmt.Sprintf(str1+"0")
		  }
		  cnt3++
		}
		v_cidr--
	}
	
	s := strings.Split(str1, ".")
	
	i := fmt.Sprintf("%0-8v", s[3])
	s3 := fmt.Sprintf(strings.ReplaceAll(i, " ", "0"))
	binarySubnetMask := fmt.Sprintf("%v.%v.%v.%v", s[0], s[1], s[2], s3)
	
	cidrToSubnetMaskMap["Binary Subnet Mask"] = binarySubnetMask
	
	n1, _ := strconv.ParseInt(s[0], 2, 64)
	n2, _ := strconv.ParseInt(s[1], 2, 64)
	n3, _ := strconv.ParseInt(s[2], 2, 64)
	n4, _ := strconv.ParseInt(s3, 2, 64)
	
	subnetMask := fmt.Sprintf("%v.%v.%v.%v", n1, n2, n3, n4)
	cidrToSubnetMaskMap["Subnet Mask"] = subnetMask
	
	// Performing bitwise XOR operation below to invert the binary digits and find the wild mask.
	// Example: 192 XOR 255 => 63    /    11000000  XOR  11111111 => 00111111
	// Refer: xor.pw
	wm1 := n1 ^ int64(maxOctetDecimal)
	wm2 := n2 ^ int64(maxOctetDecimal)
	wm3 := n3 ^ int64(maxOctetDecimal)
	wm4 := n4 ^ int64(maxOctetDecimal)
	
	wildcardMask := fmt.Sprintf("%v.%v.%v.%v", wm1, wm2, wm3, wm4)
	cidrToSubnetMaskMap["Wildcard Mask"] = wildcardMask
	
	cidrToSubnetMaskMap["Network Bits (total masked bits)"] = fmt.Sprintf("%v", cidr)
	hostsBits := ipTotalBitCount - cidr
	strHostsBits := fmt.Sprintf("%v", hostsBits)
	cidrToSubnetMaskMap["Hosts Bits (unmasked bits)"] = strHostsBits
	
	return cidrToSubnetMaskMap, nil
}




