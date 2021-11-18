#!/usr/bin/env ruby

require 'yaml'
require 'pp'

def usage
  puts "Usage: #{__FILE__} </path/to/selectorsyncset.yaml>"
  puts "Writes resulting YAML files to CWD"
  puts "Format: {selector index}_{yaml index within}_Kind_{Kind counter}.yaml"
  puts "Ex: 001_001_Namespace_1.yaml"
end


yamlfile = ARGV[0]
if yamlfile.nil?
  usage
  exit 1
end


y = YAML::load_file(yamlfile)

# Each "item" will have a different clusterDeploymentSelector, so we'll bucket by them
y['items'].each_with_index do |ss_item,i|
  puts i.to_s.rjust(3,"0")
  pp ss_item['spec']['clusterDeploymentSelector']
  puts
  c = 0
  khist = Hash.new(0)
  ss_item['spec']['resources'].each do |resource|
    kind = resource['kind']
    khist[kind] += 1
    c += 1
    File.open(i.to_s.rjust(3,"0") + "_" + c.to_s.rjust(3,"0") + "_" + kind + "_" + khist[kind].to_s + ".yaml", "w") do |f|
      f.write resource.to_yaml
    end
  end
end